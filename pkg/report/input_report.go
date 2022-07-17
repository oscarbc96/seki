package report

import (
	"encoding/json"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type InputReport struct {
	Input  load.Input
	Checks []check.CheckResult
}

func (ir *InputReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Path          string              `json:"path"`
		DetectedTypes []load.DetectedType `json:"detected_types"`
		Checks        []check.CheckResult `json:"checks"`
	}{
		Path:          ir.Input.Path(),
		DetectedTypes: ir.Input.DetectedTypes(),
		Checks:        ir.Checks,
	})
}
