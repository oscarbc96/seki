package run

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type CheckOutput struct {
	Id             string
	Name           string
	Description    string
	Severity       check.Severity
	Controls       map[string][]string
	Tags           []string
	RemediationDoc string
	Status         check.Status
	Location       load.Range
	Context        map[string]string
}
type Output struct {
	Path          string
	DetectedTypes []load.DetectedType
	Checks        []CheckOutput
}
