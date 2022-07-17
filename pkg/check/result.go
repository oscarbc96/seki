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
	Err       ErrorWrapper
}

func (cr *CheckResult) MarshalJSON() ([]byte, error) {
	type aux struct {
		Id             string              `json:"id"`
		Name           string              `json:"name"`
		Description    string              `json:"description"`
		Severity       Severity            `json:"severity"`
		Controls       map[string][]string `json:"controls,omitempty"`
		Tags           []string            `json:"tags,omitempty"`
		RemediationDoc string              `json:"remediation_doc,omitempty"`
		Status         Status              `json:"status"`
		Locations      Locations           `json:"location,omitempty"`
		Context        map[string]string   `json:"context,omitempty"`
		Err            *ErrorWrapper       `json:"error,omitempty"`
	}
	tmp := aux{}
	tmp.Id = cr.Check.Id()
	tmp.Name = cr.Check.Name()
	tmp.Description = cr.Check.Description()
	tmp.Severity = cr.Check.Severity()
	tmp.Controls = cr.Check.Controls()
	tmp.Tags = cr.Check.Tags()
	tmp.RemediationDoc = cr.Check.RemediationDoc()
	tmp.Status = cr.Status
	tmp.Locations = cr.Locations
	tmp.Context = cr.Context
	if !cr.Err.IsEmpty() {
		tmp.Err = &cr.Err
	}
	return json.Marshal(tmp)
}

func NewSkipCheckResult(check Check) CheckResult {
	return CheckResult{
		Check:  check,
		Status: SKIP,
	}
}

func NewSkipCheckResultWithError(check Check, err error) CheckResult {
	return CheckResult{
		Check:  check,
		Status: SKIP,
		Err:    ErrorWrapper{Err: err},
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
