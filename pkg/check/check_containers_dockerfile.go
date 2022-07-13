package check

import (
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
)

func init() {
	allChecks = append(
		allChecks,
		new(CheckContainersDockerfileDockerHubRateLimit),
		new(CheckContainersDockerfileLatestTag),
	)
}

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
					Line:   stage.Location[0].Start.Line, // TODO validate the hardcoded 0
					Column: stage.Location[0].Start.Character,
				},
				End: load.Position{
					Line:   stage.Location[0].End.Line,
					Column: stage.Location[0].End.Character,
				},
			},
		})
	}
	return matches, nil
}

type CheckContainersDockerfileDockerHubRateLimit struct{}

func (CheckContainersDockerfileDockerHubRateLimit) Id() string { return "SK_3" }

func (CheckContainersDockerfileDockerHubRateLimit) Name() string { return "Name" }

func (CheckContainersDockerfileDockerHubRateLimit) Description() string { return "Description" }

func (CheckContainersDockerfileDockerHubRateLimit) Severity() Severity { return Medium }

func (CheckContainersDockerfileDockerHubRateLimit) Controls() map[string][]string {
	return map[string][]string{}
}

func (CheckContainersDockerfileDockerHubRateLimit) Tags() []string { return []string{} }

func (CheckContainersDockerfileDockerHubRateLimit) RemediationDoc() string { return "RemediationDoc" }

func (CheckContainersDockerfileDockerHubRateLimit) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (CheckContainersDockerfileDockerHubRateLimit) Run(f load.Input) (*Result, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipResult(), err
	}
	defer file.Close()

	dockerLayers, err := parseDockerLayers(file)
	if err != nil {
		return NewSkipResult(), err
	}

	var locations []load.Range
	for _, layer := range dockerLayers {
		if layer.Registry == "docker.io" {
			locations = append(locations, layer.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailResult(locations), nil
	}

	return NewPassResult(), nil
}

type CheckContainersDockerfileLatestTag struct{}

func (CheckContainersDockerfileLatestTag) Id() string { return "SK_4" }

func (CheckContainersDockerfileLatestTag) Name() string { return "Name" }

func (CheckContainersDockerfileLatestTag) Description() string { return "Description" }

func (CheckContainersDockerfileLatestTag) Severity() Severity { return Medium }

func (CheckContainersDockerfileLatestTag) Controls() map[string][]string {
	return map[string][]string{}
}

func (CheckContainersDockerfileLatestTag) Tags() []string { return []string{"docker"} }

func (CheckContainersDockerfileLatestTag) RemediationDoc() string { return "RemediationDoc" }

func (CheckContainersDockerfileLatestTag) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (CheckContainersDockerfileLatestTag) Run(f load.Input) (*Result, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipResult(), err
	}
	defer file.Close()

	dockerLayers, err := parseDockerLayers(file)
	if err != nil {
		return NewSkipResult(), err
	}

	var locations []load.Range
	for _, layer := range dockerLayers {
		if layer.Tag == "latest" {
			locations = append(locations, layer.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailResult(locations), nil
	}

	return NewPassResult(), nil
}
