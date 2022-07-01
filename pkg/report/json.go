package report

import (
	"encoding/json"
	"github.com/oscarbc96/seki/pkg/result"
)

func JSONFormat(results []result.CheckResult) (string, error) {
	jsonOutput, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
