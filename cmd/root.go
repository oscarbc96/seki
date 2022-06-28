package cmd

import (
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
}

func init() {
	rootCmd.PersistentFlags().String("loggingLevel", zerolog.LevelInfoValue, "set the logging level")
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
