package cmd

import (
	"fmt"
	"github.com/oscarbc96/seki/pkg/metadata"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Seki",
	Long:  `All software has versions. This is Seki's`,
	Run: func(cmd *cobra.Command, args []string) {
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(writer, "Seki - ", metadata.Homepage)
		fmt.Fprintln(writer, "Version:", "\t", metadata.Version)
		fmt.Fprintln(writer, "Git commit:", "\t", metadata.Commit)
		fmt.Fprintln(writer, "Built:", "\t", metadata.BuiltDate)
		fmt.Fprintln(writer, "OS / Arch:", "\t", metadata.RuntimeOS, "/", metadata.RuntimeArch)
		writer.Flush()
	},
}
