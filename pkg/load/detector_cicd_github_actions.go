package load

import (
	"github.com/nektos/act/pkg/model"
	"github.com/rs/zerolog/log"
	"regexp"
)

var GitHubWorkflowsPathRegex *regexp.Regexp

func init() {
	allDetectors = append(
		allDetectors,
		new(DetectorCICDGitHubActions),
	)
	GitHubWorkflowsPathRegex = regexp.MustCompile(`[\w-/ ]*.github\/workflows\/[\w-]+\.ya?ml$`)
}

type DetectorCICDGitHubActions struct{}

func (DetectorCICDGitHubActions) Detect(input Input) (DetectedType, error) {
	log.Debug().Str("path", input.Path()).Msg("Detecting CICD GitHub Actions")

	if input.IsDir() {
		return DetectedUnknown, nil
	}

	if !GitHubWorkflowsPathRegex.MatchString(input.Path()) {
		return DetectedUnknown, nil
	}

	reader, err := input.Open()
	if err != nil {
		return DetectedUnknown, err
	}
	defer reader.Close()

	workflow, err := model.ReadWorkflow(reader)
	if err != nil {
		return DetectedUnknown, err
	}

	if len(workflow.On()) == 0 {
		return DetectedUnknown, nil
	}

	return DetectedCICDGitHubActions, nil
}
