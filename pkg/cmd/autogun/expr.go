package autogun

import (
	"time"

	"github.com/Carbonfrost/autogun/pkg/automation"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
	"github.com/chromedp/chromedp"
)

func Exprs() []*cli.Expr {
	return []*cli.Expr{
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
	}
}

func RunSource(source string) cli.Evaluator {
	return wrapTaskAsEvaluator(runSource(source))
}

func Flow(name string) cli.Evaluator {
	return wrapTaskAsEvaluator(flow(name))
}

func NavigateForward() cli.Evaluator {
	return wrapTaskAsEvaluator(chromedp.NavigateForward())
}

func NavigateBack() cli.Evaluator {
	return wrapTaskAsEvaluator(chromedp.NavigateBack())
}

func Sleep(d time.Duration) cli.Evaluator {
	return wrapTaskAsEvaluator(chromedp.Sleep(d))
}

func Reload() cli.Evaluator {
	return wrapTaskAsEvaluator(chromedp.Reload())
}

func Stop() cli.Evaluator {
	return wrapTaskAsEvaluator(chromedp.Stop())
}

func ensurePrinter(e *cli.Expression) *cli.Expression {
	// TODO In the future, printing output from the workflow is implied behavior
	return e
}

func wrapTaskAsEvaluator(act chromedp.Action) cli.EvaluatorFunc {
	return func(_ *cli.Context, v any, yield func(any) error) error {
		appendTask(v.(*automation.Automation), act)
		return yield(v)
	}
}

func appendTask(a *automation.Automation, t chromedp.Action) {
	a.Tasks = append(a.Tasks, t)
}
