package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"text/tabwriter"
	"time"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var Version = "development"
var Commit = "development"
var BuiltDate = time.Now().UTC().Format(time.RFC3339)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Seki",
	Long:  `All software has versions. This is Seki's`,
	Run: func(cmd *cobra.Command, args []string) {
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(writer, "Version:", "\t", Version)
		fmt.Fprintln(writer, "Git commit:", "\t", Commit)
		fmt.Fprintln(writer, "Built:", "\t", BuiltDate)
		fmt.Fprintln(writer, "OS / Arch:", "\t", runtime.GOOS, "/", runtime.GOARCH)
		writer.Flush()
	},
}
