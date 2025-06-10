// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun

import (
	"cmp"
	"strings"
	"time"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
	"github.com/Carbonfrost/joe-cli/extensions/structure"
)

type ScreenshotArgs struct {
	File          string            `mapstructure:"file"`
	Scale         *float64          `mapstructure:"scale"`
	Selector      string            `mapstructure:"selector"`
	By            config.SelectorBy `mapstructure:"by"`
	On            config.SelectorOn `mapstructure:"on"`
	RetryInterval *time.Duration    `mapstructure:"retry_interval"`
	AtLeast       *int              `mapstructure:"at_least"`
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
			Evaluate: bind.Evaluator(RunSource, bind.String("file")),
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
			Evaluate: bind.Evaluator(Eval, bind.String("script")),
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
			Evaluate: bind.Evaluator(Navigate, bind.String("url")),
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
			Evaluate: bind.Evaluator(Flow, bind.String("name")),
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
			Evaluate: bind.Evaluator(Sleep, bind.Duration("duration")),
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
			Evaluate: bind.Evaluator(Screenshot, bind.Value[*ScreenshotArgs]("options")),
		},
		{
			Name:     "title", // -title
			HelpText: "store the title of the current page",
			Evaluate: Title("title"),
		},
	}
}

func Screenshot(s *ScreenshotArgs) expr.Evaluator {
	scale := 1.0
	if s.Scale != nil {
		scale = *s.Scale
	}

	// TODO This selector logic is likely temporary until a selector expression is finalized
	var selectors []*config.Selector
	if s.Selector != "" {
		selectors = []*config.Selector{
			{
				Target: s.Selector,
				By:     s.By,
				On:     s.On,
			},
		}
	}

	return wrapDeferredTaskAsEvaluator(&config.Screenshot{
		Name:      cmp.Or(s.File, "screenshot.png"),
		Scale:     scale,
		Selectors: selectors,
		Options: &config.Options{
			RetryInterval: s.RetryInterval,
			AtLeast:       s.AtLeast,
		},
	})
}

func RunSource(source string) expr.Evaluator {
	return wrapTaskAsEvaluator(runSource(source))
}

func Navigate(url string) expr.Evaluator {
	nav, _ := navigate(url)
	// TODO Handle this error
	return wrapTaskAsEvaluator(nav)
}

func Flow(name string) expr.Evaluator {
	return wrapTaskAsEvaluator(flow(name))
}

func Eval(script string) expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.Eval{
		Script: script,
		Name:   "_1",
	})
}

func NavigateForward() expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.NavigateForward{})
}

func NavigateBack() expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.NavigateBack{})
}

func Sleep(d time.Duration) expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.Sleep{Duration: d})
}

func Reload() expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.Reload{})
}

func Stop() expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.Stop{})
}

func Title(name string) expr.Evaluator {
	return wrapDeferredTaskAsEvaluator(&config.Title{Name: name})
}

func ensurePrinter(e *expr.Expression) *expr.Expression {
	// TODO In the future, printing output from the workflow is implied behavior
	return e
}

func wrapDeferredTaskAsEvaluator(act config.Task) expr.EvaluatorFunc {
	return func(_ *cli.Context, v any, yield func(any) error) error {
		a := v.(*automation.Automation)
		task, err := deferredTask(act)
		if err != nil {
			return err
		}
		appendTask(a, task)
		return yield(v)
	}
}

func wrapTaskAsEvaluator(act automation.Task) expr.EvaluatorFunc {
	return func(_ *cli.Context, v any, yield func(any) error) error {
		appendTask(v.(*automation.Automation), act)
		return yield(v)
	}
}

func appendTask(a *automation.Automation, t automation.Task) {
	a.Tasks = append(a.Tasks, t)
}

func deferredTask(act config.Task) (automation.Task, error) {
	// TODO Should obtain the appropriate binder
	return automation.UsingChromedp.BindTask(act)
}
