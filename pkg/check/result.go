package check

import "github.com/oscarbc96/seki/pkg/load"

type CheckResult struct {
	Status   Status
	Location []load.Range
	Context  map[string]string
}

func NewSkipCheckResult() CheckResult {
	return CheckResult{
		Status: SKIP,
	}
}

func NewPassCheckResult() CheckResult {
	return CheckResult{
		Status: PASS,
	}
}

func NewFailCheckResult(locations []load.Range) CheckResult {
	return CheckResult{
		Status:   FAIL,
		Location: locations,
	}
}
