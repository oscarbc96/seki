package cmd

import (
	"fmt"
	check_docker_dockerfile "github.com/oscarbc96/seki/pkg/check/docker/dockerfile"
	"github.com/oscarbc96/seki/pkg/report"
	"github.com/spf13/cobra"
)

var analyseCmd = &cobra.Command{
	Use: "analyse",
	Run: func(cmd *cobra.Command, args []string) {
		formatRawValue, err := cmd.Flags().GetString("format")
		CheckErr(err)
		format, err := report.FormatFromString(formatRawValue)
		CheckErr(err)
		formater, err := report.GetFormater(format)

		result, err := check_docker_dockerfile.CheckRegistryIsAllowed()
		CheckErr(err)

		output, err := formater(result)
		CheckErr(err)
		fmt.Print(output)
	},
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.Flags().String("format", report.DefaultFormat, "set the output format")
}
