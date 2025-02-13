package autogun

import (
	"github.com/Carbonfrost/autogun/pkg/contextual"
	"github.com/Carbonfrost/autogun/pkg/workspace"
	cli "github.com/Carbonfrost/joe-cli"
)

// SetupWorkspace provides the --workspace-dir flag and ensures that prior
// to the command running, the workspace is present in the context
func SetupWorkspace() cli.Action {
	return &cli.Setup{
		Before: func(c *cli.Context) error {
			ws := &workspace.Workspace{
				Directory: ".",
			}
			c.SetContext(contextual.With(c.Context(), ws))
			return nil
		},
	}
}
