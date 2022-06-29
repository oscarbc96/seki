package result

type Severity int

const (
	Unknown Severity = iota
	Informational
	Low
	Medium
	High
	Critical
	Off
)

func (s Severity) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Informational:
		return "Informational"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	case Off:
		return "Off"
	}
	return "N/A"
}

type Result int

const (
	PASS Result = iota
	FAIL
)

type CheckResult struct {
	//Controls           []string               `json:"controls"`
	//Families           []string               `json:"families"`
	Filepath  string `json:"filepath"`
	InputType string `json:"input_type"`
	Provider  string `json:"provider"`
	//ResourceID         string                 `json:"resource_id"`
	//ResourceType       string                 `json:"resource_type"`
	//ResourceTags       map[string]interface{} `json:"resource_tags"`
	//RuleDescription    string                 `json:"rule_description"`
	//RuleID             string                 `json:"rule_id"`
	Message string `json:"message"`
	//RuleName           string                 `json:"rule_name"`
	//RuleRawResult      bool                   `json:"rule_raw_result"`
	//RuleRemediationDoc string                 `json:"rule_remediation_doc,omitempty"`
	//Result         string                 `json:"rule_result"`
	Severity Severity `json:"severity"`
	Result   Result   `json:"result"`
	//RuleSummary        string                 `json:"rule_summary"`
}
