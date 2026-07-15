// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/model"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var urlPrefix = []string{
	"http://",
	"https://",
	"about:",
	"chrome:",
}

type AutomationQuery struct {
	Automation *model.Automation

	// Selectors is the current selector set, updated by the -select expression
	// and propagated to each subsequent task that targets elements.
	Selectors []*model.Selector
}

type RunParams struct {
	Expression      *expr.Expression
	AutomationQuery *AutomationQuery
}

func Run() cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "run",
			HelpText: "Run the specified automations",
		},
		bind.Call2(runSpec, bind.Context(), useRunParams()),
	)
}

func useRunParams() bind.ActionBinder[*RunParams] {
	return bind.NewActionBinder(
		cli.Pipeline(
			FlagsAndArgs(),
			cli.AddArgs([]*cli.Arg{
				{
					Name:  "sources",
					Value: cli.List(),
					NArg:  cli.TakeUntilNextFlag,
				},
				{
					Name: "expression",
					Value: &expr.Expression{
						Exprs: Exprs(),
					},
					Options: cli.SortedExprs,
				},
			}...),
		),
		bind.Func[*RunParams](func(c *cli.Context) (*RunParams, error) {
			// Build a single automation given the source identified in the context,
			// which are loaded as files (or navigation URLs). The expression evaluation
			// functions append additional tasks.
			auto, err := convertSources(c)
			if err != nil {
				return nil, err
			}

			return &RunParams{
				Expression: ensurePrinter(expr.FromContext(c, "expression")),
				AutomationQuery: &AutomationQuery{
					Automation: auto,
				},
			}, nil
		}),
	)
}

func runSpec(ctx *cli.Context, c *RunParams) error {
	exp := c.Expression
	err := exp.Evaluate(ctx, c.AutomationQuery)
	if err != nil {
		return err
	}

	ws := FromContext(ctx)
	mo, err := ws.Load()
	if err != nil {
		return err
	}

	driver, err := automation.Bind(
		mo,
		automation.WithProtocol(automation.ProtocolChromedp),
		automation.WithAllocator(ws.EnsureAllocator()),
	)
	if err != nil {
		return err
	}

	results, err := driver.Execute(ctx, c.AutomationQuery.Automation)
	if err != nil {
		return err
	}

	// Generate output from the output variables and persistent files
	data, _ := json.MarshalIndent(results.Outputs, "", "    ")
	os.Stdout.Write(data)
	results.PersistOutputFiles()

	return nil
}

func convertSources(c *cli.Context) (*model.Automation, error) {
	var result []model.Task
	for _, source := range c.List("sources") {
		switch {
		case source == ".":
			continue

		case source == "-":
			// TODO Support reading the automation from an input file
			return nil, fmt.Errorf("not yet implemented: read automation from stdin")

		case looksLikeURL(source):
			nav, err := navigate(source)
			if err != nil {
				return nil, err
			}
			result = append(result, nav)

		default:
			result = append(result, &model.Source{Filename: source})
		}
	}
	return &model.Automation{
		Tasks: result,
	}, nil
}

func navigate(u string) (model.Task, error) {
	urlExp, _ := parseHCL(u)
	return &model.Navigate{
		URL: model.ExpressionFromHCL(urlExp),
	}, nil
}

func parseHCL(u string) (hcl.Expression, error) {
	return hclsyntax.ParseExpression([]byte(strconv.Quote(u)), "-", hcl.Pos{})
}

func looksLikeURL(addr string) bool {
	for _, prefix := range urlPrefix {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}
