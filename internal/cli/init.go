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
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 1 {
			return nil
		}
		return fmt.Errorf("expected one argument only (module)")
	},
	RunE: func(_ *cobra.Command, args []string) error {
		module := args[0]

		components := []component.Component{
			// Catalystgo config file
			component.NewConfigComponent(),

			// Go
			component.NewGomodComponent(module, strings.TrimPrefix(runtime.Version(), "go")),
			component.NewGitignoreComponent(),

			// Buf
			component.NewBufComponent(),
			component.NewBufGenComponent(),
			component.NewBufGenVendorComponent(),

			// Docker
			component.NewDockerComponent(module),
			component.NewDockerComposeComponent(module),

			// PreCommit
			component.NewPreCommitComponent(),
			component.NewcommitLintConfigComponent(),

			// Lint
			component.NewReviveComponent(),

			// Automation
			component.NewTaskfileComponent(module),
		}

		srv.Init(components, override)
		return nil
	},
}
