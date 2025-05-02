// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package automation

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type contextKey string

const (
	evalContextKey      contextKey = "evalContext"
	automationResultKey contextKey = "automationResult"
)

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

func mustAutomationResult(c context.Context) *Result {
	return c.Value(automationResultKey).(*Result)
}

func withAutomationResult(c context.Context, ar *Result) context.Context {
	return context.WithValue(c, automationResultKey, ar)
}

func withEvalContext(c context.Context) context.Context {
	return context.WithValue(c, evalContextKey, &hcl.EvalContext{})
}
