// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// Expression is an indirection over a configuration expression.
type Expression interface {
	// Value evaluates the expression within the given scope, producing a
	// value. A nil scope is treated as an empty scope.
	Value(*Scope) (cty.Value, error)

	// Variables enumerates the names of the variables that the expression
	// references.
	Variables() []string
}

// Scope provides the variables and functions available when evaluating an
// [Expression].
type Scope struct {
	Variables map[string]cty.Value
	Functions map[string]function.Function
}

// ExpressionFromHCL wraps an HCL expression as a model [Expression]. It returns
// nil when the given expression is nil.
func ExpressionFromHCL(expr hcl.Expression) Expression {
	if expr == nil {
		return nil
	}
	return hclExpression{expr: expr}
}

type hclExpression struct {
	expr hcl.Expression
}

func (e hclExpression) Value(scope *Scope) (cty.Value, error) {
	v, diags := e.expr.Value(scope.evalContext())
	if diags.HasErrors() {
		return v, diags
	}
	return v, nil
}

func (e hclExpression) Variables() []string {
	traversals := e.expr.Variables()
	names := make([]string, 0, len(traversals))
	for _, t := range traversals {
		names = append(names, t.RootName())
	}
	return names
}

func (s *Scope) evalContext() *hcl.EvalContext {
	if s == nil {
		return &hcl.EvalContext{}
	}
	return &hcl.EvalContext{
		Variables: s.Variables,
		Functions: s.Functions,
	}
}
