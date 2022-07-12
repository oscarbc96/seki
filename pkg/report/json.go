package report

import (
	"encoding/json"
	"github.com/oscarbc96/seki/pkg/run"
)

func JSONFormat(results []run.Output) (string, error) {
	jsonOutput, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
