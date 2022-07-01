package check_docker_dockerfile

import (
	"encoding/json"
	"fmt"
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/result"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"regexp"
)

func init() {
	check.Analysers["CheckRegistryIsAllowed"] = CheckRegistryIsAllowed
}

var (
	FS          = afero.NewOsFs()
	FSUtil      = &afero.Afero{Fs: FS}
	namePattern = regexp.MustCompile(`(?m)[\w.]*Dockerfile[\w.]*`)
)

type DockerLayer struct {
	Digest   string `json:"digest"`
	Image    string `json:"image"`
	Platform string `json:"platform"`
	Registry string `json:"registry"`
	Tag      string `json:"tag"`
}

func CheckRegistryIsAllowed() (*result.CheckResult, error) {
	file, err := FSUtil.Open("Dockerfile")
	if err != nil {
		return nil, err
	}

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

	var matches []DockerLayer
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

		matches = append(matches, DockerLayer{
			Digest:   digest,
			Image:    reference.Path(ref),
			Platform: stage.Platform,
			Registry: reference.Domain(ref),
			Tag:      ref.(reference.Tagged).Tag(),
		})
	}
	jsonOutput, err := json.MarshalIndent(matches, "", "  ")
	fmt.Println(string(jsonOutput))
	return nil, nil

}
