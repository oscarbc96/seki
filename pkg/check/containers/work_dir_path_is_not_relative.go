package containers

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/utils"
	"path/filepath"
)

type WorkDirPathIsNotRelative struct{}

func (WorkDirPathIsNotRelative) Id() string { return "SK_7" }

func (WorkDirPathIsNotRelative) Name() string {
	return "WORKDIR path must be absolute"
}

func (WorkDirPathIsNotRelative) Description() string { return "Description" }

func (WorkDirPathIsNotRelative) Severity() check.Severity { return check.Medium }

func (WorkDirPathIsNotRelative) Controls() map[string][]string {
	return map[string][]string{}
}

func (WorkDirPathIsNotRelative) Tags() []string { return []string{"docker"} }

func (WorkDirPathIsNotRelative) RemediationDoc() string { return "https://sekisecurity.com/docs/" }

func (WorkDirPathIsNotRelative) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c WorkDirPathIsNotRelative) Run(f load.Input) check.CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if workDirCommand, isWorkDirCommand := command.(*instructions.WorkdirCommand); isWorkDirCommand {
				if !filepath.IsAbs(workDirCommand.Path) {
					cmdLocation := command.Location()[0] // TODO validate the hardcoded 0
					colStart, colEnd := utils.FindStartAndEndColumn(workDirCommand.String(), workDirCommand.Path)
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
	}

	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
