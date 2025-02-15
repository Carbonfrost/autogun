package autogun

import (
	"github.com/Carbonfrost/autogun/pkg/automation"
	cli "github.com/Carbonfrost/joe-cli"
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
			Evaluate: bindString("file", RunSource),
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
			Evaluate: bindString("name", Flow),
		},
		{
			Name:     "forward", // -forward
			HelpText: "navigate forward in history",
			Evaluate: NavigateForward(),
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

func bindString(arg string, fn func(string) cli.Evaluator) cli.EvaluatorFunc {
	return func(c *cli.Context, v interface{}, yield func(interface{}) error) error {
		return fn(c.String(arg)).Evaluate(c, v, yield)
	}
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
