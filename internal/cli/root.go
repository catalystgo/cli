package cli

import (
	"context"

	"github.com/catalystgo/cli/internal/service"
	"github.com/catalystgo/logger/log"
	"github.com/spf13/cobra"
)

// Common flags
var (
	override bool
	verbose  bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&override, "override", "r", false, "override existing files")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	rootCmd.AddCommand(initCmd)
}

var srv = service.New()

var rootCmd = &cobra.Command{
	Use:   "catalystgo [command]",
	Short: "Catalystgo framework CLI tool for code generation",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if verbose {
			log.SetLevel(log.LevelDebug)
		}
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
