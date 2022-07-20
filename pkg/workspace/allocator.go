package workspace

import (
	"context"

	"github.com/chromedp/chromedp"
)

type contextKey string

type Allocator struct {
	URL string
}

func allocatorFromContext(c context.Context) *Allocator {
	return FromContext(c).ensureAllocator()
}

func (a *Allocator) newContext(parent context.Context) (context.Context, context.CancelFunc) {
	if a.URL != "" {
		allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(parent, a.URL)
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

func (a *Allocator) SetURL(s string) error {
	a.URL = s
	return nil
}
