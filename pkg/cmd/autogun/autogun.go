package autogun

import (
	"github.com/Carbonfrost/autogun/pkg/internal/build"
	cli "github.com/Carbonfrost/joe-cli"
)

func Run(args []string) {
	NewApp().Run(args)
}

func NewApp() *cli.App {
	return &cli.App{
		Name:     "autogun",
		HelpText: "Use Autogun to detect and execute Web browser automation",
		Comment:  "Web browser automation",
		Uses:     cli.Pipeline(),
		Action: func(c *cli.Context) error {
			c.Stdout.WriteString("Hello, world!")
			return nil
		},
		Flags: []*cli.Flag{
			{
				Name:     "chdir",
				HelpText: "Change directory into the specified working {DIRECTORY}",
				Value:    cli.String(),
				Options:  cli.WorkingDirectory | cli.NonPersistent,
			},
		},
		Version: build.Version,
	}
}
