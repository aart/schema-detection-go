package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "schema-detector",
	Short: "A CLI tool to detect schemas from various data sources.",
	Long: `A command-line interface (CLI) tool for inferring data schemas,
primarily focusing on JSON-like structures and generating BigQuery-compatible
schema definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is given
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be available to all subcommands in this application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.schema-detector.yaml)")
}
