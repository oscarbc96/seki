package report

import (
	"encoding/json"
	"github.com/oscarbc96/seki/pkg/result"
)

func JSONFormat(r *result.RuleResult) (string, error) {
	jsonOutput, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
