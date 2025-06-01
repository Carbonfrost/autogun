// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun

import (
	"context"
	"fmt"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/contextual"
	cli "github.com/Carbonfrost/joe-cli"
)

func FlagsAndArgs() cli.Action {
	return cli.Pipeline(
		cli.AddFlags([]*cli.Flag{
			{Uses: SetBrowserURL()},
			{Uses: SetEngine()},
			{Uses: SetDeviceID()},
		}...),
	)
}

func ListDevices() cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "list-devices",
			HelpText: "list available devices to emulate",
			Options:  cli.Exits | cli.NonPersistent,
			Value:    new(bool),
		},
		cli.At(cli.ActionTiming, cli.ActionOf(func() {
			for _, s := range automation.DeviceIDs() {
				fmt.Println(s)
			}
		})))
}

func SetBrowserURL(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "browser",
			Aliases:  []string{"b"},
			HelpText: "connect to the running browser instance by {URL}",
		},
		withBinding((*automation.Allocator).SetBrowserURL, v...),
	)
}

func SetEngine(v ...Engine) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "engine",
			Aliases:  []string{"e"},
			HelpText: "use the specified {ENGINE} (chromedp)",
		},
		withBinding(setEngine, v...),
	)
}

func SetDeviceID(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "device",
			Aliases:  []string{"D"},
			HelpText: "use the specified device {ID}",
		},
		withBinding((*automation.Allocator).SetDeviceID, v...),
	)
}

func withBinding[V any](binder func(*automation.Allocator, V) error, args ...V) cli.Action {
	return cli.BindContext(allocatorFromContext, binder, args...)
}

func allocatorFromContext(c context.Context) *automation.Allocator {
	return contextual.Workspace(c).EnsureAllocator()
}

func setEngine(a *automation.Allocator, e Engine) error {
	return a.SetEngine(e.Value())
}
