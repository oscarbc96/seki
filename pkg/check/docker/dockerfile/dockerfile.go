package check_docker_dockerfile

import (
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/result"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
)

type dockerLayer struct {
	Digest   string
	Image    string
	Platform string
	Registry string
	Tag      string
	Location load.Range
}

func parseDockerLayers(file io.Reader) ([]dockerLayer, error) {
	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	stages, _, err := instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return nil, err
	}

	layerNames := lo.Map[instructions.Stage, string](stages, func(stage instructions.Stage, _ int) string {
		return stage.Name
	})

	var matches []dockerLayer
	for _, stage := range stages {
		// Ignoring if stage inherits from a previous stage
		if lo.Contains[string](layerNames, stage.BaseName) {
			continue
		}

		ref, err := reference.ParseNormalizedNamed(stage.BaseName)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", stage.BaseName)
		}
		ref = reference.TagNameOnly(ref)

		var digest string
		if canonicalReference, isCanonical := ref.(reference.Canonical); isCanonical {
			digest = canonicalReference.Digest().String()
		}

		matches = append(matches, dockerLayer{
			Digest:   digest,
			Image:    reference.Path(ref),
			Platform: stage.Platform,
			Registry: reference.Domain(ref),
			Tag:      ref.(reference.Tagged).Tag(),
			Location: load.Range{
				Start: load.Position{
					Line:      stage.Location[0].Start.Line,
					Character: stage.Location[0].Start.Character,
				},
				End: load.Position{
					Line:      stage.Location[0].End.Line,
					Character: stage.Location[0].End.Character,
				},
			},
		})
	}
	return matches, nil
}

func CheckDockerHubRateLimit(f load.InputFile) ([]result.CheckResult, error) {
	file, err := f.Open()
	defer file.Close()
	if err != nil {
		return []result.CheckResult{}, err
	}

	dockerLayers, err := parseDockerLayers(file)
	if err != nil {
		return nil, err
	}

	var output []result.CheckResult
	for _, layer := range dockerLayers {
		if layer.Registry == "docker.io" {
			output = append(
				output,
				result.CheckResult{
					Result:   result.FAIL,
					Severity: result.Low,
					Message:  "Docker Hub may apply rate limiting",
					Range:    layer.Location,
				},
			)
		}
	}

	return output, nil
}

func CheckLatestTag(f load.InputFile) ([]result.CheckResult, error) {
	file, err := f.Open()
	defer file.Close()
	if err != nil {
		return []result.CheckResult{}, err
	}

	dockerLayers, err := parseDockerLayers(file)
	if err != nil {
		return nil, err
	}

	var output []result.CheckResult
	for _, layer := range dockerLayers {
		if layer.Tag == "latest" {
			output = append(
				output,
				result.CheckResult{
					Result:   result.FAIL,
					Severity: result.High,
					Message:  "Ensure the image uses a non latest version tag",
					Range:    layer.Location,
				},
			)
		}
	}

	return output, nil
}
