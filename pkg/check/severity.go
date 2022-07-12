package check

import "encoding/json"

type Severity int

const (
	Unknown Severity = iota
	Informational
	Low
	Medium
	High
	Critical
)

var severityToString = map[Severity]string{
	Unknown:       "UNKNOWN",
	Informational: "INFORMATIONAL",
	Low:           "LOW",
	Medium:        "MEDIUM",
	High:          "HIGH",
	Critical:      "CRITICAL",
}

func (s Severity) String() string {
	return severityToString[s]
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
