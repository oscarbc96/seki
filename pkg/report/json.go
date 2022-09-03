package report

import (
	"encoding/json"
)

func JSONFormatter(inputReports []InputReport) (string, error) {
	jsonOutput, err := json.MarshalIndent(inputReports, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
