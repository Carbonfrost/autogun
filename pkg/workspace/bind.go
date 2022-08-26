package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/chromedp/chromedp"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
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
		return chromedp.ActionFunc(func(c context.Context) error {
			v, err := evalContext(c, t.URL)
			if err != nil {
				return nil
			}
			fmt.Printf("Navigate to %s\n", v.AsString())
			return chromedp.Navigate(v.AsString()).Do(c)
		})
	case *config.WaitVisible:
		return bindSelector(chromedp.WaitVisible, t.Selector, t.Selectors)
	case *config.Click:
		return bindSelector(chromedp.Click, t.Selector, t.Selectors)
	case *config.Eval:
		return chromedp.ActionFunc(func(c context.Context) error {
			var msg json.RawMessage
			res.Outputs[t.Name] = &msg
			err := chromedp.Evaluate(t.Script, &msg).Do(c)
			fmt.Printf("Evaluate %v\n", string(msg))
			evalContextFrom(c).Variables[t.Name] = umarshalData(msg)
			return err
		})

	default:
		panic(fmt.Errorf("unexpected task type %T", t))
	}
}

func umarshalData(msg json.RawMessage) cty.Value {
	// TODO Additional types might occur
	m := map[string]string{}
	_ = json.Unmarshal(msg, &m)
	v, _ := gocty.ToCtyValue(m, cty.Map(cty.String))
	return v
}

func bindSelector(fn func(interface{}, ...chromedp.QueryOption) chromedp.QueryAction, sel string, sels []*config.Selector) chromedp.Tasks {
	if sel != "" {
		sels = append(sels, &config.Selector{
			Target: sel,
			By:     config.BySearch,
		})
	}

	tasks := make([]chromedp.Action, len(sels))
	for i, s := range sels {
		opts := make([]chromedp.QueryOption, 0)
		if s.By != "" {
			opts = append(opts, bindSelectorBy(s.By))
		}
		if s.On != "" {
			opts = append(opts, bindSelectorOn(s.On))
		}

		tasks[i] = fn(s.Target, opts...)
	}
	return tasks
}

func bindSelectorBy(s config.SelectorBy) chromedp.QueryOption {
	switch s {
	case config.BySearch:
		return chromedp.BySearch
	case config.ByJSPath:
		return chromedp.ByJSPath
	case config.ByID:
		return chromedp.ByID
	case config.ByQuery:
		return chromedp.ByQuery
	case config.ByQueryAll:
		return chromedp.ByQueryAll
	}
	return nil
}

func bindSelectorOn(s config.SelectorOn) chromedp.QueryOption {
	switch s {
	case config.OnReady:
		return chromedp.NodeReady
	case config.OnVisible:
		return chromedp.NodeVisible
	case config.OnNotVisible:
		return chromedp.NodeNotVisible
	case config.OnEnabled:
		return chromedp.NodeEnabled
	case config.OnSelected:
		return chromedp.NodeSelected
	case config.OnNotPresent:
		return chromedp.NodeNotPresent
	}
	return nil
}

func (r *AutomationResult) bindAutomation(automation *config.Automation) []chromedp.Action {
	actions := make([]chromedp.Action, 0)
	for _, t := range automation.Tasks {
		actions = append(actions, bindTask(r, t))
	}
	return actions
}
