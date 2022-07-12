package load

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/rs/zerolog/log"
	"strings"
)

func init() {
	allDetectors = append(
		allDetectors,
		new(DetectorContainersDockerfile),
	)
}

type DetectorContainersDockerfile struct{}

func (DetectorContainersDockerfile) Detect(input Input) (DetectedType, error) {
	log.Debug().Str("path", input.Path()).Msg("Detecting Containers Dockerfile")

	if input.IsDir() {
		return DetectedUnknown, nil
	}

	if !strings.Contains(strings.ToLower(input.Name()), "dockerfile") {
		return DetectedUnknown, nil
	}

	file, err := input.Open()
	if err != nil {
		return DetectedUnknown, err
	}
	defer file.Close()
	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return DetectedUnknown, err
	}

	_, _, err = instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return DetectedUnknown, err
	}
	return DetectedContainerDockerfile, nil
}
