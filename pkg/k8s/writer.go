package k8s

import (
	"fmt"

	"github.com/aquasecurity/trivy/pkg/k8s/report"

	cdx "github.com/CycloneDX/cyclonedx-go"

	rp "github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/report/table"
)

type Writer interface {
	Write(report.Report) error
}

// Write writes the results in the give format
func Write(k8sreport report.Report, option report.Option) error {
	k8sreport.PrintErrors()

	switch option.Format {
	case rp.FormatJSON:
		jwriter := report.JSONWriter{
			Output: option.Output,
			Report: option.Report,
		}
		return jwriter.Write(k8sreport)
	case rp.FormatTable:
		separatedReports := report.SeparateMisconfigReports(k8sreport, option.Scanners, option.Components)

		if option.Report == report.SummaryReport {
			target := fmt.Sprintf("Summary Report for %s", k8sreport.ClusterName)
			table.RenderTarget(option.Output, target, table.IsOutputToTerminal(option.Output))
		}

		for _, r := range separatedReports {
			writer := &report.TableWriter{
				Output:        option.Output,
				Report:        option.Report,
				Severities:    option.Severities,
				ColumnHeading: report.ColumnHeading(option.Scanners, option.Components, r.Columns),
			}

			if err := writer.Write(r.Report); err != nil {
				return err
			}
		}

		return nil
	case rp.FormatCycloneDX:
		w := report.NewCycloneDXWriter(option.Output, cdx.BOMFileFormatJSON, option.APIVersion)
		return w.Write(k8sreport.RootComponent)
	}
	return nil
}
