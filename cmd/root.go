package cmd

import (
	"fmt"
	check_aws_cloudformation_s3 "github.com/oscarbc96/seki/pkg/check/aws/cloudformation/s3"
	check_docker_dockerfile "github.com/oscarbc96/seki/pkg/check/docker/dockerfile"
	"github.com/oscarbc96/seki/pkg/report"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "seki",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

		loggingLevelRawValue, err := cmd.Flags().GetString("loggingLevel")
		CheckErr(err)
		loggingLevel, err := zerolog.ParseLevel(loggingLevelRawValue)
		CheckErr(err)
		zerolog.SetGlobalLevel(loggingLevel)
	},

	Run: func(cmd *cobra.Command, args []string) {
		formatRawValue, err := cmd.Flags().GetString("format")
		CheckErr(err)
		format, err := report.FormatFromString(formatRawValue)
		CheckErr(err)
		formater, err := report.GetFormater(format)

		result, err := check_aws_cloudformation_s3.CheckS3ObjectVersioningRule()
		CheckErr(err)

		output, err := formater(result)
		CheckErr(err)
		fmt.Print(output)

		result, err = check_docker_dockerfile.CheckRegistryIsAllowed()
		CheckErr(err)

		output, err = formater(result)
		CheckErr(err)
		fmt.Print(output)
	},
}

func init() {
	rootCmd.PersistentFlags().String("loggingLevel", zerolog.LevelInfoValue, "set the logging level")
	rootCmd.Flags().String("format", report.DefaultFormat, "set the output format")
}

func CheckErr(msg interface{}) {
	if msg != nil {
		log.Error().Msgf("%s", msg)
		os.Exit(1)
	}
}

func Execute() {
	rootCmd.Execute()
}
