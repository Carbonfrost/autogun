package workspace

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/chromedp/chromedp"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type AutomationResult struct {
	Outputs     map[string]*json.RawMessage
	OutputFiles map[string]*[]byte
}

type produceQueryActionFunc = func(interface{}, ...chromedp.QueryOption) chromedp.QueryAction
type produceFileUserActionFunc = func(*[]byte) chromedp.Action

func NewAutomationResult() *AutomationResult {
	return &AutomationResult{
		Outputs:     map[string]*json.RawMessage{},
		OutputFiles: map[string]*[]byte{},
	}
}

func bindTask(task config.Task) chromedp.Action {
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
		return bindSelector(chromedp.WaitVisible, t.Selector, t.Selectors, t.Options)
	case *config.Click:
		return bindSelector(chromedp.Click, t.Selector, t.Selectors, t.Options)
	case *config.Screenshot:
		filename := cmp.Or(t.Name, "screenshot.png")
		if t.Selector == "" && len(t.Selectors) == 0 {
			return requestOutputFile(filename, chromedp.CaptureScreenshot)
		}

		return bindSelector(func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
			curry := func(file *[]byte) chromedp.Action {
				return chromedp.Screenshot(sel, file, opts...)
			}
			return requestOutputFile(filename, curry)
		}, t.Selector, t.Selectors, t.Options)

	case *config.Eval:
		return chromedp.ActionFunc(func(c context.Context) error {
			res := mustAutomationResult(c)
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

func bindSelector(fn produceQueryActionFunc, sel string, sels []*config.Selector, options *config.Options) chromedp.Tasks {
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
		opts = append(opts, bindQueryOptions(options)...)

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

func bindQueryOptions(opts *config.Options) (results []chromedp.QueryOption) {
	if opts == nil {
		return
	}
	if opts.RetryInterval != nil {
		results = append(results, chromedp.RetryInterval(*opts.RetryInterval))
	}
	if opts.AtLeast != nil {
		results = append(results, chromedp.AtLeast(*opts.AtLeast))
	}
	return
}

func bindAutomation(automation *config.Automation) []chromedp.Action {
	actions := make([]chromedp.Action, 0)
	for _, t := range automation.Tasks {
		actions = append(actions, bindTask(t))
	}
	return actions
}

// TODO This should not necessarily be API
func (r *AutomationResult) PersistOutputFiles() {
	for name, f := range r.OutputFiles {
		file, err := os.Create(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error persisting output files: %v\n", err)
			continue
		}

		_, err = io.Copy(file, bytes.NewReader(*f))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error persisting output files: %v\n", err)
			continue
		}
	}
}

func requestOutputFile(name string, fn produceFileUserActionFunc) chromedp.Action {
	return chromedp.ActionFunc(func(c context.Context) error {
		r := mustAutomationResult(c)
		file := make([]byte, 2048)
		r.OutputFiles[name] = &file

		act := fn(&file)
		return act.Do(c)
	})
}

func mustAutomationResult(c context.Context) *AutomationResult {
	return c.Value(automationResultKey).(*AutomationResult)
}

func withAutomationResult(c context.Context, ar *AutomationResult) context.Context {
	return context.WithValue(c, automationResultKey, ar)
}
