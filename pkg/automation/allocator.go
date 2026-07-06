// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Carbonfrost/autogun/pkg/model"
)

type Allocator struct {
	BrowserURL string
	Engine     Protocol
	DeviceID   string

	// Options carries the union of exec/remote allocator options. Fields
	// that do not pertain to the selected allocator produce a warning on
	// stderr when the context is created.
	Options *AllocatorOptions
}

func (a *Allocator) SetBrowserURL(s string) error {
	a.BrowserURL = s
	return nil
}

func (a *Allocator) SetEngine(e Protocol) error {
	a.Engine = e
	return nil
}

func (a *Allocator) SetDeviceID(v string) error {
	a.DeviceID = v

	if v != "" {
		_, ok := a.resolveDevice()
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: device %q not found\n", a.DeviceID)
		}
	}
	return nil
}

// ensureOptions lazily initializes the allocator options.
func (a *Allocator) ensureOptions() *AllocatorOptions {
	if a.Options == nil {
		a.Options = &AllocatorOptions{}
	}
	return a.Options
}

// SetExecPath sets the browser executable path (exec allocator).
func (a *Allocator) SetExecPath(v string) error {
	a.ensureOptions().ExecPath = &v
	return nil
}

// SetProxyServer sets the outbound proxy server (exec allocator).
func (a *Allocator) SetProxyServer(v string) error {
	a.ensureOptions().ProxyServer = &v
	return nil
}

// SetUserAgent sets the default User-Agent header (exec allocator).
func (a *Allocator) SetUserAgent(v string) error {
	a.ensureOptions().UserAgent = &v
	return nil
}

// SetUserDataDir sets the Chrome profile directory (exec allocator).
func (a *Allocator) SetUserDataDir(v string) error {
	a.ensureOptions().UserDataDir = &v
	return nil
}

// SetHeadless toggles headless mode (exec allocator).
func (a *Allocator) SetHeadless(v bool) error {
	a.ensureOptions().Headless = &v
	return nil
}

// SetNoSandbox toggles disabling the sandbox (exec allocator).
func (a *Allocator) SetNoSandbox(v bool) error {
	a.ensureOptions().NoSandbox = &v
	return nil
}

// SetDisableGPU toggles disabling the GPU process (exec allocator).
func (a *Allocator) SetDisableGPU(v bool) error {
	a.ensureOptions().DisableGPU = &v
	return nil
}

// SetIgnoreCertErrors toggles ignoring certificate errors (exec allocator).
func (a *Allocator) SetIgnoreCertErrors(v bool) error {
	a.ensureOptions().IgnoreCertErrors = &v
	return nil
}

// SetWSURLReadTimeout sets the websocket URL read timeout (exec allocator).
func (a *Allocator) SetWSURLReadTimeout(v time.Duration) error {
	a.ensureOptions().WSURLReadTimeout = &v
	return nil
}

// SetNoModifyURL toggles preventing the remote allocator from rewriting the
// websocket URL (remote allocator).
func (a *Allocator) SetNoModifyURL(v bool) error {
	a.ensureOptions().NoModifyURL = &v
	return nil
}

// SetEnv appends NAME=value environment variables passed to the browser
// process (exec allocator).
func (a *Allocator) SetEnv(v []string) error {
	o := a.ensureOptions()
	o.Env = append(o.Env, v...)
	return nil
}

// SetWindowSize sets the initial window size, given as "WxH" or "W,H"
// (exec allocator).
func (a *Allocator) SetWindowSize(v string) error {
	if v == "" {
		return nil
	}
	sep := strings.IndexAny(v, "x,")
	if sep < 0 {
		return fmt.Errorf("invalid window size %q: expected WxH", v)
	}
	width, err := strconv.Atoi(strings.TrimSpace(v[:sep]))
	if err != nil {
		return fmt.Errorf("invalid window size %q: %w", v, err)
	}
	height, err := strconv.Atoi(strings.TrimSpace(v[sep+1:]))
	if err != nil {
		return fmt.Errorf("invalid window size %q: %w", v, err)
	}
	a.ensureOptions().WindowSize = &WindowSize{Width: width, Height: height}
	return nil
}

// SetFlag adds arbitrary Chrome command-line flags, each given as "name" (a
// boolean flag) or "name=value" (exec allocator).
func (a *Allocator) SetFlag(v []string) error {
	o := a.ensureOptions()
	if o.Flags == nil {
		o.Flags = make(map[string]any)
	}
	for _, entry := range v {
		name, value, ok := strings.Cut(entry, "=")
		if ok {
			o.Flags[name] = value
		} else {
			o.Flags[name] = true
		}
	}
	return nil
}

func (a *Allocator) newContext(parent context.Context) (context.Context, context.CancelFunc, error) {
	ctx := withEvalContext(parent)
	eng := a.Engine
	if eng == nil {
		eng = ProtocolChromedp
	}

	if a.BrowserURL != "" {
		return eng.NewRemoteAllocator(ctx, a.BrowserURL, a.Options)
	}

	return eng.NewExecAllocator(ctx, a.Options)
}

func (a *Allocator) resolveDevice() (dev model.Device, ok bool) {
	dev, ok = model.LookupDevice(a.DeviceID)
	return
}
