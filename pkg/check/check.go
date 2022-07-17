package check

import (
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/samber/lo"
)

type Check interface {
	Id() string
	Name() string
	Description() string
	Severity() Severity
	Controls() map[string][]string
	Tags() []string
	RemediationDoc() string
	InputTypes() []load.DetectedType
	Run(f load.Input) (CheckResult, error)
}

var allChecks = []Check{}

func GetChecksFor(tpe load.DetectedType) []Check {
	return lo.Filter[Check](allChecks, func(check Check, _ int) bool {
		return lo.Contains[load.DetectedType](check.InputTypes(), tpe)
	})
}
