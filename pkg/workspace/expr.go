// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"cmp"
	"strings"
	"time"

	internalcli "github.com/Carbonfrost/autogun/pkg/internal/cli"
	"github.com/Carbonfrost/autogun/pkg/model"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
	"github.com/Carbonfrost/joe-cli/extensions/structure"
)

type ScreenshotArgs struct {
	File          string           `mapstructure:"file"`
	Scale         *float64         `mapstructure:"scale"`
	Selector      string           `mapstructure:"selector"`
	By            model.SelectorBy `mapstructure:"by"`
	On            model.SelectorOn `mapstructure:"on"`
	RetryInterval *time.Duration   `mapstructure:"retry_interval"`
	AtLeast       *int             `mapstructure:"at_least"`
}

func Exprs() []*expr.Expr {
	return []*expr.Expr{
		{
			Name:     "run", // -run FILE
			HelpText: "run an automation from a FILE",
			Args: []*cli.Arg{
				{
					Name:  "file",
					Value: new(string), // TODO This should be *cli.File
					NArg:  1,
				},
			},
			Evaluate: expr.BindEvaluator(RunSource, bind.String("file")),
		},
		{
			Name:     "eval", // -eval SCRIPT
			HelpText: "evaluate a script",
			Args: []*cli.Arg{
				{
					Name:      "script",
					Value:     new(string),
					NArg:      1,
					UsageText: "SCRIPT | @FILE",
					Options:   cli.AllowFileReference,
				},
			},
			Evaluate: expr.BindEvaluator(Eval, bind.String("script")),
		},
		{
			Name:     "navigate", // -navigate URL
			HelpText: "navigate to the specified {URL}",
			Args: []*cli.Arg{
				{
					Name:  "url",
					Value: new(string),
					NArg:  1,
				},
			},
			Evaluate: expr.BindEvaluator(Navigate, bind.String("url")),
		},
		{
			Name:     "flow", // -flow NAME
			HelpText: "run an automation by NAME",
			Args: []*cli.Arg{
				{
					Name:  "name",
					Value: new(string),
					NArg:  1,
				},
			},
			Evaluate: expr.BindEvaluator(Flow, bind.String("name")),
		},
		{
			Name:     "forward", // -forward
			HelpText: "navigate forward in history",
			Evaluate: NavigateForward(),
		},
		{
			Name:     "back", // -back
			HelpText: "navigate back in history",
			Evaluate: NavigateBack(),
		},
		{
			Name:     "reload", // -reload
			HelpText: "reload the current page",
			Evaluate: Reload(),
		},
		{
			Name:     "stop", // -stop
			HelpText: "stop loading the current page",
			Evaluate: Stop(),
		},
		{
			Name:     "sleep", // -sleep DURATION
			HelpText: "sleep for the DURATION",
			Args: []*cli.Arg{
				{
					Name:  "duration",
					Value: new(time.Duration),
					NArg:  1,
				},
			},
			Evaluate: expr.BindEvaluator(Sleep, bind.Duration("duration")),
		},
		{
			Name:     "select", // -select SELECTORS
			HelpText: "set the selector set applied to subsequent operations",
			Args: []*cli.Arg{
				{
					Name:      "selectors",
					Value:     new(internalcli.SelectorSet),
					NArg:      1,
					UsageText: "SELECTORS",
				},
			},
			Evaluate: expr.BindEvaluator(Select, bind.Value[[]*model.Selector]("selectors")),
		},
		{
			Name:     "blur", // -blur
			HelpText: "blur the selected element",
			Evaluate: Blur(),
		},
		{
			Name:     "clear", // -clear
			HelpText: "clear the selected element",
			Evaluate: Clear(),
		},
		{
			Name:     "click", // -click
			HelpText: "click the selected element",
			Evaluate: Click(),
		},
		{
			Name:     "doubleclick", // -doubleclick
			HelpText: "double-click the selected element",
			Evaluate: DoubleClick(),
		},
		{
			Name:     "send_keys", // -send_keys KEYS
			HelpText: "send {KEYS} to the selected element",
			Args: []*cli.Arg{
				{
					Name:  "keys",
					Value: new(string),
					NArg:  1,
				},
			},
			Evaluate: expr.BindEvaluator(SendKeys, bind.String("keys")),
		},
		{
			Name:     "wait_visible", // -wait_visible
			HelpText: "wait for the selected element to become visible",
			Evaluate: WaitVisible(),
		},
		{
			Name:     "screenshot", // -screenshot [scale=SCALE,]
			HelpText: "capture a screenshot",
			Args: []*cli.Arg{
				{
					Name:  "options",
					Value: structure.Of(new(ScreenshotArgs)),
					NArg: cli.OptionalArg(func(s string) bool {
						return !strings.HasPrefix(s, "-")
					}),
				},
			},
			Evaluate: expr.BindEvaluator(Screenshot, bind.Value[*ScreenshotArgs]("options")),
		},
		{
			Name:     "title", // -title
			HelpText: "store the title of the current page",
			Evaluate: Title("title"),
		},
		{
			Name:     "inner_html", // -inner_html
			HelpText: "store the inner HTML of the selected element",
			Evaluate: InnerHTML("inner_html"),
		},
		{
			Name:     "version", // -version
			HelpText: "print out version information",
			Evaluate: Version(),
		},
		{
			Name:     "options", // -options retry_interval=TIME,at_least=1
			Aliases:  []string{"O"},
			HelpText: "specify options for handling selector",
			Args: []*cli.Arg{
				{
					Name:      "value",
					Value:     structure.Of(new(model.Options)),
					UsageText: "{retry_interval=TIME,at_least=NUM}",
				},
			},
			Evaluate: expr.BindEvaluator(SetOptions, bind.Value[model.Options]("value")),
		},
	}
}

func SetOptions(m model.Options) expr.Evaluator {
	return withQuery(func(q *AutomationQuery) error {
		q.Options = &m
		return nil
	})
}

func Screenshot(s *ScreenshotArgs) expr.Evaluator {
	scale := 1.0
	if s.Scale != nil {
		scale = *s.Scale
	}

	// The local selector, when set, takes precedence over the selector set from
	// the AutomationQuery.
	var local []*model.Selector
	if s.Selector != "" {
		local = []*model.Selector{
			{
				Target: s.Selector,
				By:     s.By,
				On:     s.On,
			},
		}
	}

	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		if local != nil {
			selectors = local
		}

		var localOpts model.Options
		if opts != nil {
			localOpts = *opts
			if s.RetryInterval != nil {
				localOpts.RetryInterval = s.RetryInterval
			}
			if s.AtLeast != nil {
				localOpts.AtLeast = s.AtLeast
			}
		}
		return &model.Screenshot{
			Name:      cmp.Or(s.File, "screenshot.png"),
			Scale:     scale,
			Selectors: selectors,
			Options:   &localOpts,
		}
	})
}

func Select(set []*model.Selector) expr.Evaluator {
	return withQuery(func(q *AutomationQuery) error {
		q.Selectors = set
		return nil
	})
}

func Blur() expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.Blur{Selectors: selectors, Options: opts}
	})
}

func Clear() expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.Clear{Selectors: selectors, Options: opts}
	})
}

func Click() expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.Click{Selectors: selectors, Options: opts}
	})
}

func DoubleClick() expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.DoubleClick{Selectors: selectors, Options: opts}
	})
}

func SendKeys(keys string) expr.Evaluator {
	keysExp, _ := parseHCL(keys)
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.SendKeys{
			Selectors: selectors,
			Options:   opts,
			Keys:      model.ExpressionFromHCL(keysExp),
		}
	})
}

func WaitVisible() expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.WaitVisible{Selectors: selectors, Options: opts}
	})
}

func RunSource(source string) expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Source{Filename: source})
}

func Navigate(url string) expr.Evaluator {
	nav, _ := navigate(url)
	// TODO Handle this error
	return wrapTaskAsEvaluator(nav)
}

func Flow(name string) expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Flow{Name: name})
}

func Eval(script string) expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Eval{
		Script: script,
		Name:   "_1",
	})
}

func NavigateForward() expr.Evaluator {
	return wrapTaskAsEvaluator(&model.NavigateForward{})
}

func NavigateBack() expr.Evaluator {
	return wrapTaskAsEvaluator(&model.NavigateBack{})
}

func Sleep(d time.Duration) expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Sleep{Duration: d})
}

func Reload() expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Reload{})
}

func Stop() expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Stop{})
}

func Title(name string) expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Title{Name: name})
}

func InnerHTML(name string) expr.Evaluator {
	return wrapSelectorTask(func(selectors []*model.Selector, opts *model.Options) model.Task {
		return &model.InnerHTML{Name: name, Selectors: selectors, Options: opts}
	})
}

func Version() expr.Evaluator {
	return wrapTaskAsEvaluator(&model.Version{})
}

func ensurePrinter(e *expr.Expression) *expr.Expression {
	// TODO In the future, printing output from the workflow is implied behavior
	return e
}

func wrapTaskAsEvaluator(act model.Task) expr.EvaluatorFunc {
	return withAutomation(func(a *model.Automation) error {
		appendTask(a, act)
		return nil
	})
}

func withAutomation(fn func(*model.Automation) error) expr.EvaluatorFunc {
	return withQuery(func(query *AutomationQuery) error {
		return fn(query.Automation)
	})
}

func withQuery(fn func(*AutomationQuery) error) expr.EvaluatorFunc {
	return func(_ *cli.Context, v any, yield func(any) error) error {
		query := v.(*AutomationQuery)
		if err := fn(query); err != nil {
			return err
		}
		return yield(query)
	}
}

func wrapSelectorTask(build func(selectors []*model.Selector, opts *model.Options) model.Task) expr.EvaluatorFunc {
	return withQuery(func(query *AutomationQuery) error {
		appendTask(query.Automation, build(query.Selectors, query.Options))
		return nil
	})
}

func appendTask(a *model.Automation, t model.Task) {
	a.Tasks = append(a.Tasks, t)
}
