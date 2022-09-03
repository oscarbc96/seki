package containers

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/utils"
	"github.com/samber/lo"
)

type LastUserIsNotRoot struct{}

func (LastUserIsNotRoot) Id() string { return "SK_6" }

func (LastUserIsNotRoot) Name() string {
	return "Ensure the last USER is not root"
}

func (LastUserIsNotRoot) Description() string { return "Description" }

func (LastUserIsNotRoot) Severity() check.Severity { return check.Medium }

func (LastUserIsNotRoot) Controls() map[string][]string {
	return map[string][]string{}
}

func (LastUserIsNotRoot) Tags() []string { return []string{"docker"} }

func (LastUserIsNotRoot) RemediationDoc() string { return "https://sekisecurity.com/docs/" }

func (LastUserIsNotRoot) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c LastUserIsNotRoot) Run(f load.Input) check.CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	lastStage := stages[len(stages)-1]

	userCommands := lo.FilterMap[instructions.Command, instructions.UserCommand](lastStage.Commands, func(command instructions.Command, _ int) (instructions.UserCommand, bool) {
		if userCommand, isUserCommand := command.(*instructions.UserCommand); isUserCommand {
			return *userCommand, true
		}
		return instructions.UserCommand{}, false
	})

	if len(userCommands) == 0 {
		return check.NewPassCheckResult(c)
	}

	lastUserCommand := userCommands[len(userCommands)-1]
	if lastUserCommand.User == "root" {
		colStart, colEnd := utils.FindStartAndEndColumn(lastUserCommand.String(), "root")
		cmdLocation := lastUserCommand.Location()[0] // TODO validate the hardcoded 0

		return check.NewFailCheckResult(c, []load.Range{
			{
				Start: load.Position{
					Line:   cmdLocation.Start.Line,
					Column: colStart,
				},
				End: load.Position{
					Line:   cmdLocation.End.Line,
					Column: colEnd,
				},
			},
		})
	}

	return check.NewPassCheckResult(c)
}
