package autogun

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Carbonfrost/autogun/pkg/contextual"
	"github.com/Carbonfrost/autogun/pkg/workspace"
	cli "github.com/Carbonfrost/joe-cli"
)

var (
	expectedOneArg = errors.New("expected zero or one arguments")
)

func FlagsAndArgs() cli.Action {
	return cli.Pipeline(
		cli.AddFlags([]*cli.Flag{
			{Uses: SetBrowserURL()},
		}...),
	)

}

func RunAutomation() cli.Action {
	return cli.Setup{
		Action: func(c *cli.Context) error {
			ws := contextual.Workspace(c)
			err := ws.Load(c.FileSet("files").Files...)
			if err != nil {
				return err
			}

			for _, auto := range ws.Automations() {
				res, err := ws.ExecuteCore(auto)
				if err != nil {
					return err
				}

				data, _ := json.MarshalIndent(res.Outputs, "", "    ")
				os.Stdout.Write(data)
			}
			return nil
		},
	}
}

func SetBrowserURL(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "browser",
			Aliases:  []string{"b"},
			HelpText: "Connect to the running browser instance by {URL}",
		},
		withBinding((*workspace.Allocator).SetBrowserURL, v),
	)
}

func withBinding[V any](binder func(*workspace.Allocator, V) error, args []V) cli.Action {
	switch len(args) {
	case 0:
		return cli.BindContext(allocatorFromContext, binder)
	case 1:
		return cli.BindContext(allocatorFromContext, binder, args[0])
	default:
		panic(expectedOneArg)
	}
}

func allocatorFromContext(c context.Context) *workspace.Allocator {
	return contextual.Workspace(c).EnsureAllocator()
}
