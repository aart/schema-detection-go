package cmd

import (
	"fmt"
	"os"
	"strings"

	"schema-detection-go/core"

	"github.com/spf13/cobra"
)

var (
	numWorkers int
)

func init() {
	runCmd.Flags().IntVarP(&numWorkers, "workers", "w", 4, "Number of workers to use for concurrent processing")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [files...]",
	Short: "Run schema detection on specified files",
	Long: `This command processes one or more files to infer their schema.
Each file is expected to contain JSON lines.`,
	Args: cobra.MinimumNArgs(1), // Requires at least one file argument
	Run: func(cmd *cobra.Command, args []string) {
		if numWorkers <= 0 {
			fmt.Println("Number of workers must be greater than 0")
			os.Exit(1)
		}

		fmt.Printf("Processing files: %s with %d workers...\n", strings.Join(args, ", "), numWorkers)

		// Read lines from all specified files concurrently
		linesChannel := core.ReadLines(args, numWorkers)

		// Infer schema concurrently from the lines
		finalSchema := core.InferSchemaConcurrently(linesChannel, numWorkers)

		if finalSchema == nil {
			fmt.Println("Could not infer a schema from the provided files.")
			os.Exit(1)
		}

		// Print the resulting schema (for now, a simple representation)
		fmt.Println("Inferred Schema:")
		for _, field := range finalSchema.Fields {
			fmt.Printf("  - Name: %s, Type: %s, Repeated: %t, Required: %t\n",
				field.Name, field.Type, field.Repeated, field.Required)
			if field.Type == "RECORD" && len(field.Fields) > 0 {
				printSubFields(field.Fields, "    ")
			}
		}
	},
}

func printSubFields(fields []*core.FieldSchema, indent string) {
	for _, field := range fields {
		fmt.Printf("%s- Name: %s, Type: %s, Repeated: %t, Required: %t\n",
			indent, field.Name, field.Type, field.Repeated, field.Required)
		if field.Type == "RECORD" && len(field.Fields) > 0 {
			printSubFields(field.Fields, indent+"  ")
		}
	}
}
