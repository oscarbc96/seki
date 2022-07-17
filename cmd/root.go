package cmd

import (
	"fmt"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/report"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "seki [flags] [...path]",
	Args: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			exists, err := load.PathExists(path)
			if err != nil {
				log.Fatal().Err(err)
			}
			if !exists {
				log.Fatal().Str("path", path).Msg("Path does not exist")
			}
		}
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loggingLevelRawValue, err := cmd.Flags().GetString("loggingLevel")
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get loggingLevel flag")
		}
		loggingLevel, err := zerolog.ParseLevel(loggingLevelRawValue)
		if err != nil {
			log.Fatal().Str("loggingLevelRaw", loggingLevelRawValue).Err(err).Msg("Could not parse logging level")
		}
		zerolog.SetGlobalLevel(loggingLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		formatId, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get format flag")
		}

		formatter, err := report.FormatterFromString(formatId)
		if err != nil {
			log.Fatal().Str("format", formatId).Err(err).Msg("Could not get formatter")
		}

		paths := args
		if len(paths) == 0 {
			log.Info().Msg("Path not specified, using current working directory")
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal().Err(err).Msg("Could not get current working directory")
			}
			paths = append(paths, cwd)
		}
		inputs, _ := load.FlatPathsToInputs(paths)

		var reports []report.InputReport
		for _, input := range inputs {
			inputReport := report.InputReport{Input: input}
			for _, detectedType := range input.DetectedTypes() {
				for _, chck := range check.GetChecksFor(detectedType) {
					result := chck.Run(input)
					inputReport.Checks = append(inputReport.Checks, result)
				}
			}
			if len(inputReport.Checks) > 0 {
				reports = append(reports, inputReport)
			}
		}

		output, err := formatter(reports)
		if err != nil {
			log.Error().Str("format", formatId).Err(err).Msg("Could not format")
		}

		fmt.Print(output)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("loggingLevel", "l", zerolog.LevelWarnValue, "set the logging level")
	rootCmd.Flags().StringP("format", "f", report.DefaultFormat, "set the output format")
}

func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err)
	}
}
