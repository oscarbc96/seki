package load

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/rs/zerolog/log"
	"strings"
)

func DetectDockerfile(f InputFile) (bool, error) {
	log.Debug().Str("path", f.path).Msg("Detecting Dockerfile")

	if f.IsDir() {
		return false, nil
	}

	if !strings.Contains(strings.ToLower(f.Name()), "dockerfile") {
		return false, nil
	}

	file, err := f.Open()
	defer file.Close()
	if err != nil {
		return false, err
	}
	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return false, err
	}

	_, _, err = instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return false, err
	}
	return true, nil
}
