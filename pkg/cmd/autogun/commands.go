// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autogun

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/contextual"
	internalcli "github.com/Carbonfrost/autogun/pkg/internal/cli"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
)

func FlagsAndArgs() cli.Action {
	return cli.Pipeline(
		cli.AddFlags([]*cli.Flag{
			{Uses: SetBrowserURL()},
			{Uses: SetProtocol()},
			{Uses: SetDeviceID()},
			{Uses: SetExecPath()},
			{Uses: SetProxyServer()},
			{Uses: SetUserAgent()},
			{Uses: SetUserDataDir()},
			{Uses: SetWindowSize()},
			{Uses: SetHeadless()},
			{Uses: SetNoSandbox()},
			{Uses: SetDisableGPU()},
			{Uses: SetIgnoreCertErrors()},
			{Uses: SetWSURLReadTimeout()},
			{Uses: SetEnv()},
			{Uses: SetFlag()},
			{Uses: SetNoModifyURL()},
		}...),
	)
}

func SetBrowserURL(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "browser",
			Aliases:  []string{"b"},
			HelpText: "connect to the running browser instance by {URL}",
		},
		withBinding(setBrowserURLHelper, v...),
	)
}

func setBrowserURLHelper(a *automation.Allocator, source string) error {
	// Interpreted as default protocol on localhost when just a port
	if port, ok := strings.CutPrefix(source, ":"); ok {
		source = "ws://" + net.JoinHostPort("127.0.0.1", port)
	}
	return a.SetBrowserURL(source)
}

func SetProtocol(v ...automation.Protocol) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "protocol",
			Aliases:  []string{"P"},
			HelpText: "use the specified {ENGINE} (chromedp)",
			Value:    new(internalcli.Protocol),
		},
		withBinding((*automation.Allocator).SetProtocol, v...),
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

func SetExecPath(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "exec-path",
			HelpText: "run the browser at {PATH} (exec allocator)",
		},
		withBinding((*automation.Allocator).SetExecPath, v...),
	)
}

func SetProxyServer(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "proxy-server",
			HelpText: "route browser traffic through the proxy {SERVER} (exec allocator)",
		},
		withBinding((*automation.Allocator).SetProxyServer, v...),
	)
}

func SetUserAgent(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "user-agent",
			HelpText: "set the default User-Agent header to {STRING} (exec allocator)",
		},
		withBinding((*automation.Allocator).SetUserAgent, v...),
	)
}

func SetUserDataDir(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "user-data-dir",
			HelpText: "use {DIR} as the Chrome profile directory (exec allocator)",
		},
		withBinding((*automation.Allocator).SetUserDataDir, v...),
	)
}

func SetWindowSize(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "window-size",
			HelpText: "set the initial window size, given as {WxH} (exec allocator)",
		},
		withBinding((*automation.Allocator).SetWindowSize, v...),
	)
}

func SetHeadless(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "headless",
			HelpText: "run the browser in headless mode (exec allocator)",
			Value:    new(bool),
		},
		withBinding((*automation.Allocator).SetHeadless, v...),
	)
}

func SetNoSandbox(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "no-sandbox",
			HelpText: "disable the browser sandbox (exec allocator)",
			Value:    new(bool),
		},
		withBinding((*automation.Allocator).SetNoSandbox, v...),
	)
}

func SetDisableGPU(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "disable-gpu",
			HelpText: "disable the GPU process (exec allocator)",
			Value:    new(bool),
		},
		withBinding((*automation.Allocator).SetDisableGPU, v...),
	)
}

func SetIgnoreCertErrors(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "ignore-cert-errors",
			HelpText: "ignore certificate-related errors (exec allocator)",
			Value:    new(bool),
		},
		withBinding((*automation.Allocator).SetIgnoreCertErrors, v...),
	)
}

func SetWSURLReadTimeout(v ...time.Duration) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "ws-url-timeout",
			HelpText: "wait up to {DURATION} for the websocket URL (exec allocator)",
			Value:    new(time.Duration),
		},
		withBinding((*automation.Allocator).SetWSURLReadTimeout, v...),
	)
}

func SetEnv(v ...[]string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "env",
			HelpText: "pass {NAME=VALUE} into the browser process, repeatable (exec allocator)",
			Value:    cli.List(),
		},
		withBinding((*automation.Allocator).SetEnv, v...),
	)
}

func SetFlag(v ...[]string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "flag",
			HelpText: "pass an arbitrary Chrome {NAME} or {NAME=VALUE} flag, repeatable (exec allocator)",
			Value:    cli.List(),
		},
		withBinding((*automation.Allocator).SetFlag, v...),
	)
}

func SetNoModifyURL(v ...bool) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "no-modify-url",
			HelpText: "do not rewrite the websocket URL (remote allocator)",
			Value:    new(bool),
		},
		withBinding((*automation.Allocator).SetNoModifyURL, v...),
	)
}

func withBinding[V any](binder func(*automation.Allocator, V) error, args ...V) cli.Action {
	return bind.Call2(binder, bind.FromContext(allocatorFromContext), bind.Exact(args...))
}

func allocatorFromContext(c context.Context) *automation.Allocator {
	return contextual.Workspace(c).EnsureAllocator()
}
