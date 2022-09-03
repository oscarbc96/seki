package check

import (
	"github.com/oscarbc96/seki/pkg/load"
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
	Run(f load.Input) CheckResult
}
