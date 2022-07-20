package workspace

import (
	"errors"

	cli "github.com/Carbonfrost/joe-cli"
)

var (
	expectedOneArg = errors.New("expected zero or one arguments")
)

// SetupWorkspace provides the --workspace-dir flag and ensures that prior
// to the command running, the workspace is present in the context
func SetupWorkspace() cli.Action {
	return &cli.Setup{
		Before: func(c *cli.Context) error {
			ws := &Workspace{
				Directory: ".",
			}
			c.Context = SetContextWorkspace(c.Context, ws)
			return nil
		},
	}
}

func FlagsAndArgs() cli.Action {
	return cli.Pipeline(
		cli.AddFlags([]*cli.Flag{
			{Uses: SetURL()},
		}...),
	)

}

func SetURL(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "browser",
			Aliases:  []string{"b"},
			HelpText: "Connect to the running browser instance by {URL}",
		},
		withBinding((*Allocator).SetURL, v),
	)
}

func withBinding[V any](binder func(*Allocator, V) error, args []V) cli.Action {
	switch len(args) {
	case 0:
		return cli.BindContext(allocatorFromContext, binder)
	case 1:
		return cli.BindContext(allocatorFromContext, binder, args[0])
	default:
		panic(expectedOneArg)
	}
}
