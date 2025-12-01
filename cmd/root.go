package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/owlcode3/renommer/internal/common"
	"github.com/owlcode3/renommer/internal/pathx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "renommer",
	Short: "A powerful CLI tool to rename files and directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			color.Red("Error: requires at least 2 arg(s), only received 1")
			return
		}

		if err := pathx.RoutePath(args); err != nil {
			color.Red("Error: %v", err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVarP(&common.Verbose, "verbose", "v", false, "verbose output")
}
