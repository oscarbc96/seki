package report

import (
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/oscarbc96/seki/pkg/check"
	"strings"
)

func getSeverityColor(s check.Severity) text.Color {
	switch s {
	case check.Unknown:
		return text.FgWhite
	case check.Informational:
		return text.FgHiCyan
	case check.Low:
		return text.FgHiBlue
	case check.Medium:
		return text.FgHiYellow
	case check.High:
		return text.FgHiRed
	case check.Critical:
		return text.BgHiRed
	}
	return text.Reset
}

func getStatusColor(s check.Status) text.Color {
	switch s {
	case check.PASS:
		return text.FgHiGreen
	case check.FAIL:
		return text.FgHiRed
	case check.SKIP:
		return text.FgHiBlue
	}
	return text.Reset
}

func TextFormatter(input_reports []InputReport) (string, error) {
	var result []string

	for _, input_report := range input_reports {
		for _, detectedType := range input_report.Input.DetectedTypes() {
			result = append(result, text.BgHiMagenta.Sprint(detectedType))
			result = append(result, " ")
		}
		result = append(result, text.Bold.Sprint(input_report.Input.Path()))
		result = append(result, "\n")
		for _, checkResult := range input_report.CheckResults {
			severity := checkResult.Check.Severity()
			status := checkResult.Status
			result = append(result, text.FgHiBlue.Sprint(checkResult.Check.Id()))
			result = append(result, " ")
			result = append(result, getSeverityColor(severity).Sprintf("[%s]", severity))
			result = append(result, " ")
			result = append(result, getStatusColor(status).Sprint(status))
			result = append(result, " ")
			result = append(result, text.Reset.Sprint(checkResult.Check.Name()))
			result = append(result, " ")
			result = append(result, text.FgHiCyan.Sprint(checkResult.Check.RemediationDoc()))
			result = append(result, "\n")

			for _, loc := range checkResult.Locations {
				if loc.IsEmpty() {
					continue
				}

				code, _ := input_report.Input.ReadLines(loc.Start.Line, loc.End.Line)
				for idx, line := range code {
					result = append(result, text.Bold.Sprintf("%v", loc.Start.Line+idx))
					result = append(result, " ")
					result = append(result, text.BgHiBlack.Sprint(line))
					result = append(result, "\n")
				}
			}

			if !checkResult.Err.IsEmpty() {
				result = append(result, text.FgHiRed.Sprint("Error from the check:\n"))
				result = append(result, text.Reset.Sprintf("%s\n", checkResult.Err.Err))
			}
		}
		result = append(result, "\n")
	}

	return strings.Join(result, ""), nil
}
