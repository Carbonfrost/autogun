package workspace

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type contextKey string

const evalContextKey contextKey = "evalContext"

type Allocator struct {
	BrowserURL string
}

func allocatorFromContext(c context.Context) *Allocator {
	return FromContext(c).ensureAllocator()
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

func (a *Allocator) SetBrowserURL(s string) error {
	a.BrowserURL = s
	return nil
}

func evalContext(c context.Context, expr hcl.Expression) (cty.Value, error) {
	ec := evalContextFrom(c)
	v, _ := expr.Value(ec)
	return v, nil
}

func evalContextFrom(c context.Context) *hcl.EvalContext {
	ec := c.Value(evalContextKey).(*hcl.EvalContext)
	if ec.Variables == nil {
		ec.Variables = map[string]cty.Value{}

	}
	return ec
}
