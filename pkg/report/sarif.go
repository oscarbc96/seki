package report

import (
	"bytes"
	"github.com/oscarbc96/seki/pkg/metadata"
	"github.com/oscarbc96/seki/pkg/result"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

func SARIFReport(results []result.CheckResult) (string, error) {
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return "", err
	}

	driver := sarif.NewDriver("seki").
		WithSemanticVersion(metadata.Version).
		WithInformationURI(metadata.Homepage)
	tool := sarif.NewTool(driver)
	run := sarif.NewRun(*tool)

	for _, r := range results {
		run.AddRule(r.ID).
			WithDescription(r.Description).
			WithHelpURI(r.RemediationDoc)

		run.CreateResultForRule(r.ID).
			//WithLevel(strings.ToLower(r.Severity)).
			WithMessage(sarif.NewTextMessage(r.Message)).
			WithLocations(
				[]*sarif.Location{
					sarif.NewLocationWithPhysicalLocation(
						sarif.NewPhysicalLocation().
							WithArtifactLocation(sarif.NewSimpleArtifactLocation(r.Path)).
							WithRegion(
								sarif.NewRegion().
									WithStartLine(r.Range.Start.Line).
									WithStartColumn(r.Range.Start.Character).
									WithEndLine(r.Range.End.Line).
									WithEndColumn(r.Range.End.Character),
							),
					),
				},
			)
	}
	report.AddRun(run)

	buffer := bytes.NewBuffer([]byte{})
	if err := report.PrettyWrite(buffer); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
