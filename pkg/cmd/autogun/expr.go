// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autogun

import (
	"cmp"
	"strings"
	"time"

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
			Name:     "version", // -version
			HelpText: "print out version information",
			Evaluate: Version(),
		},
	}
}

func Screenshot(s *ScreenshotArgs) expr.Evaluator {
	scale := 1.0
	if s.Scale != nil {
		scale = *s.Scale
	}

	// TODO This selector logic is likely temporary until a selector expression is finalized
	var selectors []*model.Selector
	if s.Selector != "" {
		selectors = []*model.Selector{
			{
				Target: s.Selector,
				By:     s.By,
				On:     s.On,
			},
		}
	}

	return wrapTaskAsEvaluator(&model.Screenshot{
		Name:      cmp.Or(s.File, "screenshot.png"),
		Scale:     scale,
		Selectors: selectors,
		Options: &model.Options{
			RetryInterval: s.RetryInterval,
			AtLeast:       s.AtLeast,
		},
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
	return func(_ *cli.Context, v any, yield func(any) error) error {
		query := v.(*AutomationQuery)
		a := query.Automation
		err := fn(a)
		if err != nil {
			return err
		}
		return yield(query)
	}
}

func appendTask(a *model.Automation, t model.Task) {
	a.Tasks = append(a.Tasks, t)
}
