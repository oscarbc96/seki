package cmd

import (
	"fmt"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/report"
	"github.com/oscarbc96/seki/pkg/result"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
)

var (
	FS     = afero.NewOsFs()
	FSUtil = &afero.Afero{Fs: FS}
)

var rootCmd = &cobra.Command{
	Use: "seki [flags] [...path]",
	Args: func(cmd *cobra.Command, args []string) error {
		for _, path := range args {
			exists, err := FSUtil.Exists(path)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("path doesn't exist: %s", path)
			}
		}
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loggingLevelRawValue, err := cmd.Flags().GetString("loggingLevel")
		CheckErr(err)
		loggingLevel, err := zerolog.ParseLevel(loggingLevelRawValue)
		CheckErr(err)
		zerolog.SetGlobalLevel(loggingLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		files, _ := load.ListFiles(args)

		final := []result.CheckResult{}

		for _, file := range files {
			checks := []check.CheckFunction{}

			isDocker, _ := load.DetectDockerfile(file)
			if isDocker {
				checks = append(checks, check.Checkers["dockerfile"]...)
			}
			isCloudformation, _ := load.DetectCloudformation(file)
			if isCloudformation {
				checks = append(checks, check.Checkers["cloudformation"]...)
			}

			for _, check := range checks {
				checkResult, err := check(file)
				CheckErr(err)
				final = append(final, checkResult...)
			}

		}
		formatRawValue, err := cmd.Flags().GetString("format")
		CheckErr(err)
		format, err := report.FormatFromString(formatRawValue)
		CheckErr(err)
		formater, err := report.GetFormater(format)

		output, err := formater(final)
		CheckErr(err)
		fmt.Print(output)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("loggingLevel", "l", zerolog.LevelInfoValue, "set the logging level")
	rootCmd.Flags().StringP("format", "f", report.DefaultFormat, "set the output format")
}

func CheckErr(msg interface{}) {
	if msg != nil {
		log.Error().Msgf("%s", msg)
		os.Exit(1)
	}
}

func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	CheckErr(rootCmd.Execute())
}
