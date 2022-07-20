package autogun

import (
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
		Version: "",
	}
}
