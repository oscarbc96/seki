package check

import "github.com/oscarbc96/seki/pkg/load"

type Result struct {
	Status   Status
	Location []load.Range
	Context  map[string]string
}

func NewSkipResult() *Result {
	return &Result{
		Status: SKIP,
	}
}

func NewPassResult() *Result {
	return &Result{
		Status: PASS,
	}
}

func NewFailResult(locations []load.Range) *Result {
	return &Result{
		Status:   FAIL,
		Location: locations,
	}
}
