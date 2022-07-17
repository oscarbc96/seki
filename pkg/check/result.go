package check

import (
	"encoding/json"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/samber/lo"
)

type Locations []load.Range

func (l *Locations) MarshalJSON() ([]byte, error) {
	aux := lo.Filter[load.Range](*l, func(r load.Range, _ int) bool {
		return !r.IsEmpty()
	})

	return json.Marshal(aux)
}

type CheckResult struct {
	Check     Check
	Status    Status
	Locations Locations
	Context   map[string]string
}

func (cr *CheckResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id             string              `json:"id"`
		Name           string              `json:"name"`
		Description    string              `json:"description"`
		Severity       Severity            `json:"severity"`
		Controls       map[string][]string `json:"controls,omitempty"`
		Tags           []string            `json:"tags,omitempty"`
		RemediationDoc string              `json:"remediationdoc,omitempty"`
		Status         Status              `json:"status"`
		Locations      Locations           `json:"location,omitempty"`
		Context        map[string]string   `json:"context,omitempty"`
	}{
		Id:             cr.Check.Id(),
		Name:           cr.Check.Name(),
		Description:    cr.Check.Description(),
		Severity:       cr.Check.Severity(),
		Controls:       cr.Check.Controls(),
		Tags:           cr.Check.Tags(),
		RemediationDoc: cr.Check.RemediationDoc(),
		Status:         cr.Status,
		Locations:      cr.Locations,
		Context:        cr.Context,
	})
}

func NewSkipCheckResult(check Check) CheckResult {
	return CheckResult{
		Check:  check,
		Status: SKIP,
	}
}

func NewPassCheckResult(check Check) CheckResult {
	return CheckResult{
		Check:  check,
		Status: PASS,
	}
}

func NewFailCheckResult(check Check, locations []load.Range) CheckResult {
	return CheckResult{
		Check:     check,
		Status:    FAIL,
		Locations: locations,
	}
}
