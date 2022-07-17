package report

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type InputReport struct {
	Input  load.Input
	Checks []check.CheckResult
}
