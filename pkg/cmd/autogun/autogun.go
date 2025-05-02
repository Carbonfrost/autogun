// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/internal/build"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/color"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
)

const versionTemplate = "{{ .App.Name }}, version {{ .App.Version }} {{ ExtendedVersionInfo }}\n"

var keyModules = map[string]bool{
	"cdproto":  true,
	"chromedp": true,
}

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
			versionInfoSupport(),
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
		Version: build.Version,
	}
}

func versionInfoSupport() cli.Action {
	var (
		displayExtended   bool
		extendVersionInfo = func() string {
			if !displayExtended {
				return ""
			}

			var res []string

			if info, ok := debug.ReadBuildInfo(); ok {
				for _, d := range info.Deps {
					index := strings.LastIndex(d.Path, "/")
					if index > 0 {
						name := d.Path[index+1:]
						if keyModules[name] {
							res = append(res, fmt.Sprintf("%s@%s", name, d.Version))
						}
					}
				}
			}
			return strings.Join(res, ", ")
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
	)
}
