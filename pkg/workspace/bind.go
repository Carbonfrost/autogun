package workspace

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/chromedp/chromedp"
)

type AutomationResult struct {
	Outputs map[string]*json.RawMessage
}

func NewAutomationResult() *AutomationResult {
	return &AutomationResult{
		Outputs: map[string]*json.RawMessage{},
	}
}

func RunAutomation() cli.Action {
	return cli.Setup{
		Action: func(c *cli.Context) error {
			ws := FromContext(c)
			err := ws.Load(c.FileSet("files").Files...)
			if err != nil {
				return err
			}

			for _, auto := range ws.Automations() {
				res, err := ws.executeCore(auto)
				if err != nil {
					return err
				}

				data, _ := json.MarshalIndent(res.Outputs, "", "    ")
				os.Stdout.Write(data)
			}
			return nil
		},
	}
}

func bindTask(res *AutomationResult, task config.Task) chromedp.Action {
	switch t := task.(type) {
	case *config.Navigate:
		return chromedp.Navigate(t.URL)
	case *config.Eval:
		var msg json.RawMessage
		res.Outputs[t.Name] = &msg
		return chromedp.Evaluate(t.Script, &msg)
	default:
		panic(fmt.Errorf("unexpected task type %T", t))
	}
}

func (r *AutomationResult) bindAutomation(automation *config.Automation) []chromedp.Action {
	actions := make([]chromedp.Action, 0)
	for _, t := range automation.Tasks {
		actions = append(actions, bindTask(r, t))
	}
	return actions
}
