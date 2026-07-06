// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autogun

import (
	"fmt"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/internal/build"
	"github.com/Carbonfrost/autogun/pkg/workspace"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/color"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
)

const versionTemplate = "{{ .App.Name }}, version {{ .App.Version }}{{ ExtendedVersionInfo }}\n"

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
			workspace.New(),
			versionInfoSupport(),

			// TODO It would be better to only run this implicitly when it looks
			// like a URL is being invoked
			cli.ImplicitCommand("run"),
		),
		Flags: []*cli.Flag{
			{Uses: workspace.ListDevices()},
			{Uses: SetVerbose()},
		},
		Commands: []*cli.Command{
			Fmt(),
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
						Value: &expr.Expression{
							Exprs: Exprs(),
						},
						Options: cli.SortedExprs,
					},
				},
				Uses: cli.Pipeline(
					FlagsAndArgs(),
				),
				Action: RunAutomation,
			},
			Check(),
		},
		Version: build.Version.Version,
	}
}

func SetVerbose(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "verbose",
			Aliases:  []string{"v"},
			HelpText: "Display verbose output; can be used multiple times to increase detail",
			Value:    new(bool),
		},
	)
}

func versionInfoSupport() cli.Action {
	var (
		displayExtended   bool
		extendVersionInfo = func() string {
			if !displayExtended {
				return ""
			}
			return formatVersionInfo(build.Version)
		}

		triggerExtendedVersionInfo = func(c *cli.Context) error {
			if c.Occurrences("version") >= 2 {
				displayExtended = true
			}
			return nil
		}
	)
	return cli.Pipeline(
		cli.Before(cli.ActionFunc(triggerExtendedVersionInfo)),
		cli.RegisterTemplateFunc("ExtendedVersionInfo", extendVersionInfo),
		cli.RegisterTemplate("Version", versionTemplate),
		cli.Customize("--version", cli.Alias("V")),
	)
}

func formatVersionInfo(vi build.VersionInfo) string {
	res := []string{""} // allow leading space
	for k, v := range vi.Modules {
		res = append(res, fmt.Sprintf("%v, version %v", k, v))
	}
	return strings.Join(res, "\n")
}
