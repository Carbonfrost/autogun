package automation

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type contextKey string

const evalContextKey contextKey = "evalContext"
const automationResultKey contextKey = "automationResult"

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
