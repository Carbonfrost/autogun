// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// AllocatorOptions is a data-only representation of the union of exec
// and remote allocator options which can be represented as data. Fields are
// nullable so that nil/absent means "unset" and not meant to be passed through
// to the underlying allocator.
type AllocatorOptions struct {
	// ExecPath maps to chromedp.ExecPath: the browser binary to run.
	ExecPath *string

	// Env maps to chromedp.Env: extra NAME=value environment variables.
	// nil means unset; an empty non-nil slice is treated as unset.
	Env []string

	// WSURLReadTimeout maps to chromedp.WSURLReadTimeout.
	WSURLReadTimeout *time.Duration

	// Flags is the escape hatch mapping to chromedp.Flag for arbitrary
	// Chrome command-line flags. A string value renders as --name=value;
	// a bool true renders as --name; a bool false suppresses the flag.
	// The typed convenience fields below take precedence over identical
	// keys set here.
	Flags map[string]any

	// The following are typed conveniences that chromedp exposes as
	// dedicated options but which ultimately set one or more entries in the
	// flag map.
	UserDataDir           *string
	ProxyServer           *string
	UserAgent             *string
	WindowSize            *WindowSize
	IgnoreCertErrors      *bool
	NoSandbox             *bool
	NoFirstRun            *bool
	NoDefaultBrowserCheck *bool
	Headless              *bool // chromedp.Headless (also hide-scrollbars, mute-audio)
	DisableGPU            *bool // chromedp.DisableGPU (also enable-unsafe-swiftshader)

	// Remote allocator:

	// NoModifyURL maps to chromedp.NoModifyURL. When true, the remote
	// allocator will not attempt to rewrite the provided websocket URL.
	NoModifyURL *bool
}

// WindowSize is the data form of the chromedp.WindowSize(width, height) option.
type WindowSize struct {
	Width  int
	Height int
}

func (o *AllocatorOptions) warnRemoteOnlyForExec() {
	if o.NoModifyURL != nil {
		warnf("NoModifyURL is a remote allocator option; ignored for exec allocator")
	}
}

func (o *AllocatorOptions) warnExecOnlyForRemote() {
	for _, name := range o.execOnlyFieldsSet() {
		warnf("%s is an exec allocator option; ignored for remote allocator", name)
	}
}

func (o *AllocatorOptions) execOptions() []chromedp.ExecAllocatorOption {
	var opts []chromedp.ExecAllocatorOption

	// Generic flags first so typed conveniences below can override them.
	for name, value := range o.Flags {
		opts = append(opts, chromedp.Flag(name, value))
	}

	if o.ExecPath != nil {
		opts = append(opts, chromedp.ExecPath(*o.ExecPath))
	}
	if len(o.Env) > 0 {
		opts = append(opts, chromedp.Env(o.Env...))
	}
	if o.WSURLReadTimeout != nil {
		opts = append(opts, chromedp.WSURLReadTimeout(*o.WSURLReadTimeout))
	}
	if o.UserDataDir != nil {
		opts = append(opts, chromedp.UserDataDir(*o.UserDataDir))
	}
	if o.ProxyServer != nil {
		opts = append(opts, chromedp.ProxyServer(*o.ProxyServer))
	}
	if o.UserAgent != nil {
		opts = append(opts, chromedp.UserAgent(*o.UserAgent))
	}
	if o.WindowSize != nil {
		opts = append(opts, chromedp.WindowSize(o.WindowSize.Width, o.WindowSize.Height))
	}
	if boolValue(o.IgnoreCertErrors) {
		opts = append(opts, chromedp.IgnoreCertErrors)
	}
	if boolValue(o.NoSandbox) {
		opts = append(opts, chromedp.NoSandbox)
	}
	if boolValue(o.NoFirstRun) {
		opts = append(opts, chromedp.NoFirstRun)
	}
	if boolValue(o.NoDefaultBrowserCheck) {
		opts = append(opts, chromedp.NoDefaultBrowserCheck)
	}
	if boolValue(o.Headless) {
		opts = append(opts, chromedp.Headless)
	}
	if boolValue(o.DisableGPU) {
		opts = append(opts, chromedp.DisableGPU)
	}

	return opts
}

func (o *AllocatorOptions) remoteOptions() []chromedp.RemoteAllocatorOption {
	var opts []chromedp.RemoteAllocatorOption
	if boolValue(o.NoModifyURL) {
		opts = append(opts, chromedp.NoModifyURL)
	}
	return opts
}

func (o *AllocatorOptions) execOnlyFieldsSet() []string {
	var names []string
	add := func(set bool, name string) {
		if set {
			names = append(names, name)
		}
	}
	add(o.ExecPath != nil, "ExecPath")
	add(len(o.Env) > 0, "Env")
	add(o.WSURLReadTimeout != nil, "WSURLReadTimeout")
	add(len(o.Flags) > 0, "Flags")
	add(o.UserDataDir != nil, "UserDataDir")
	add(o.ProxyServer != nil, "ProxyServer")
	add(o.UserAgent != nil, "UserAgent")
	add(o.WindowSize != nil, "WindowSize")
	add(o.IgnoreCertErrors != nil, "IgnoreCertErrors")
	add(o.NoSandbox != nil, "NoSandbox")
	add(o.NoFirstRun != nil, "NoFirstRun")
	add(o.NoDefaultBrowserCheck != nil, "NoDefaultBrowserCheck")
	add(o.Headless != nil, "Headless")
	add(o.DisableGPU != nil, "DisableGPU")
	return names
}

func boolValue(b *bool) bool {
	return b != nil && *b
}

func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}
