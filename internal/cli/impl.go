package cli

import (
	"github.com/catalystgo/logger/log"
	"github.com/spf13/cobra"
)

var (
	implementCmd = &cobra.Command{
		Use:     "implement",
		Short:   "Generate basic proto service implementation",
		Aliases: []string{"impl"},
		Run: func(cmd *cobra.Command, _ []string) {
			var (
				input  = cmd.Flag("input").Value.String()
				output = cmd.Flag("output").Value.String()
			)

			err := srv.Implement(input, output)
			if err != nil {
				log.Fatalf("implement service => %v", err)
			}
		},
	}
)

func init() {
	implementCmd.Flags().StringP("input", "i", "pkg", "input directory to generated Go files from protoc-gen-catalystgo")
	implementCmd.Flags().StringP("output", "o", "internal/api/", "output directory to save the generated Go server implementation")
	rootCmd.AddCommand(implementCmd)
}
