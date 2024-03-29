package containers

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/metadata"
	"github.com/oscarbc96/seki/utils"
)

type PreferCopyOverAdd struct{}

func (PreferCopyOverAdd) Id() string { return "SK_5" }

func (PreferCopyOverAdd) Name() string {
	return "Ensure that COPY is used instead of ADD"
}

func (PreferCopyOverAdd) Description() string { return "Description" }

func (PreferCopyOverAdd) Severity() check.Severity { return check.Medium }

func (PreferCopyOverAdd) Controls() map[string][]string {
	return map[string][]string{}
}

func (PreferCopyOverAdd) Tags() []string { return []string{"docker"} }

func (PreferCopyOverAdd) RemediationDoc() string {
	return metadata.GenerateChecksDocsURL("containers/prefer-copy-over-add")
}

func (PreferCopyOverAdd) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c PreferCopyOverAdd) Run(f load.Input) check.CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if addCommand, isAddCommand := command.(*instructions.AddCommand); isAddCommand {
				cmdLocation := command.Location()[0] // TODO validate the hardcoded 0
				colStart, colEnd := utils.FindStartAndEndColumn(addCommand.String(), "ADD")
				locations = append(locations, load.Range{
					Start: load.Position{
						Line:   cmdLocation.Start.Line,
						Column: colStart,
					},
					End: load.Position{
						Line:   cmdLocation.End.Line,
						Column: colEnd,
					},
				})
			}
		}
	}

	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
