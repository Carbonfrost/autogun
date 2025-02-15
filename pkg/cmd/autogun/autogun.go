package autogun

import (
	"github.com/Carbonfrost/autogun/pkg/internal/build"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/color"
)

func Run(args []string) {
	NewApp().Run(args)
}

func NewApp() *cli.App {
	return &cli.App{
		Name:     "autogun",
		HelpText: "Use Autogun to detect and execute Web browser automation",
		Comment:  "Web browser automation",
		Uses: cli.Pipeline(
			cli.Sorted,
			color.Options{},
			SetupWorkspace(),
		),
		Flags: []*cli.Flag{
			{
				Name:     "chdir",
				HelpText: "Change directory into the specified working {DIRECTORY}",
				Value:    new(cli.File),
				Options:  cli.WorkingDirectory | cli.NonPersistent,
			},
			{Uses: ListDevices()},
		},
		Commands: []*cli.Command{
			{
				Name:     "run",
				HelpText: "Run the specified automations",
				Args: []*cli.Arg{
					{
						Name:  "sources",
						Value: cli.List(),
						NArg:  cli.TakeUntilNextFlag,
					},
					{
						Name: "expression",
						Value: &cli.Expression{
							Exprs: Exprs(),
						},
					},
				},
				Uses: cli.Pipeline(
					FlagsAndArgs(),
				),
				Action: RunAutomation,
			},
		},
		Version: build.Version,
	}
}
