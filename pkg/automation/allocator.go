package automation

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/chromedp"
)

type Allocator struct {
	BrowserURL string
	Engine     SupportedBinder
	DeviceID   string
}

func (a *Allocator) SetBrowserURL(s string) error {
	a.BrowserURL = s
	return nil
}

func (a *Allocator) SetEngine(e SupportedBinder) error {
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

func (a *Allocator) newContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx := withEvalContext(parent)

	if a.BrowserURL != "" {
		allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(ctx, a.BrowserURL)
		res, cancelInner := chromedp.NewContext(
			allocatorContext,
		)
		return res, func() {
			defer cancelAllocator()
			defer cancelInner()
		}
	}

	return chromedp.NewContext(
		ctx,
	)
}

func (a *Allocator) resolveDevice() (dev chromedp.Device, ok bool) {
	dev, ok = devices[a.DeviceID]
	return
}
