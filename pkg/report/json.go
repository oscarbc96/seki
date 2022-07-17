package report

import (
	"encoding/json"
)

func JSONFormatter(input_reports []InputReport) (string, error) {
	jsonOutput, err := json.MarshalIndent(input_reports, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
