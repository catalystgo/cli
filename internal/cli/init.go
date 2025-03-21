package cli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/catalystgo/cli/internal/component"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize a new project",
	Long:    "Initialize a new project",
	Example: "catalystgo init github.com/user_name/project_name",
	Aliases: []string{"i"},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 1 {
			return nil
		}
		return fmt.Errorf("expected one argument only (module)")
	},
	RunE: func(_ *cobra.Command, args []string) error {
		module := args[0]

		components := []component.Component{
			// Go
			component.NewGomodComponent(module, strings.TrimPrefix(runtime.Version(), "go")),
			component.NewGitignoreComponent(),
			component.NewGoMainComponent(module),

			// Buf
			component.NewBufComponent(),
			component.NewBufGenComponent(),
			component.NewBufGenVendorComponent(),

			// Docker
			component.NewDockerComponent(module),

			// Lint
			component.NewReviveComponent(),

			// Automation
			component.NewTaskfileComponent(module),
		}

		// Config
		components = append(components, component.NewConfigComponent(module)...)

		srv.Init(components, override)
		return nil
	},
}
