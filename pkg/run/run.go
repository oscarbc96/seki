package run

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/check/cloudformation"
	"github.com/oscarbc96/seki/pkg/check/containers"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/samber/lo"
)

var AllChecks []check.Check

func GetChecksFor(tpe load.DetectedType) []check.Check {
	return lo.Filter[check.Check](AllChecks, func(check check.Check, _ int) bool {
		return lo.Contains[load.DetectedType](check.InputTypes(), tpe)
	})
}

func init() {
	groupOfChecks := [][]check.Check{
		containers.Checks,
		cloudformation.Checks,
	}

	for _, group := range groupOfChecks {
		AllChecks = append(AllChecks, group...)
	}
}
