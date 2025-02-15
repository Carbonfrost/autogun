package automation

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/hashicorp/hcl/v2"
)

type Allocator struct {
	BrowserURL string
}

func (a *Allocator) SetBrowserURL(s string) error {
	a.BrowserURL = s
	return nil
}

func (a *Allocator) newContext(parent context.Context) (context.Context, context.CancelFunc) {
	parent = context.WithValue(parent, evalContextKey, &hcl.EvalContext{})
	if a.BrowserURL != "" {
		allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(parent, a.BrowserURL)
		res, cancelInner := chromedp.NewContext(
			allocatorContext,
		)
		return res, func() {
			defer cancelAllocator()
			defer cancelInner()
		}
	}

	return chromedp.NewContext(
		parent,
	)
}
