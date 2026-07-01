// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/chromedp"
)

type Allocator struct {
	BrowserURL string
	Engine     Protocol
	DeviceID   string
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
		dev, ok := a.resolveDevice()
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: device %q not found\n", dev)
		}
	}
	return nil
}

func (a *Allocator) newContext(parent context.Context) (context.Context, context.CancelFunc, error) {
	ctx := withEvalContext(parent)
	eng := a.Engine
	if eng == nil {
		eng = UsingChromedp
	}

	if a.BrowserURL != "" {
		return eng.NewRemoteAllocator(ctx, a.BrowserURL)
	}

	return eng.NewExecAllocator(ctx)
}

func (a *Allocator) resolveDevice() (dev chromedp.Device, ok bool) {
	dev, ok = devices[a.DeviceID]
	return
}
