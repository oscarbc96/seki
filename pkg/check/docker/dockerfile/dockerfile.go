package check_docker_dockerfile

import (
	"encoding/json"
	"fmt"
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/result"
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

type Match struct {
	Digest   string `json:"digest"`
	Image    string `json:"image"`
	Layer    string `json:"layer"`
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

	var layerNames []string
	for _, stage := range stages {
		if stage.Name != "" {
			layerNames = append(layerNames, stage.Name)
		}
	}

	var matches []Match
	for _, stage := range stages {
		// check stage inherits from a previous stage
		if lo.Contains[string](layerNames, stage.BaseName) {
			continue
		}

		ref, _ := reference.ParseAnyReference(stage.BaseName)
		named := ref.(reference.Named)

		refMatch := reference.ReferenceRegexp.FindStringSubmatch(reference.TagNameOnly(named).String())

		matches = append(matches, Match{
			Digest:   refMatch[3],
			Image:    reference.Path(named),
			Layer:    stage.Name,
			Platform: stage.Platform,
			Registry: reference.Domain(named),
			Tag:      refMatch[2],
		})
	}
	jsonOutput, err := json.MarshalIndent(matches, "", "  ")
	fmt.Println(string(jsonOutput))
	return nil, nil

}
