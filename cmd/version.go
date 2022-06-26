package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var Version = "development"
var Revision = "development"
var Date = time.Now().UTC().Format(time.RFC3339)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Seki",
	Long:  `All software has versions. This is Seki's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", Version)
		fmt.Println("Revision:", Revision)
		fmt.Println("Date:", Date)
	},
}
