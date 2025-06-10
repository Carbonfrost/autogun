// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package automation

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type produceQueryActionFunc = func(any, ...chromedp.QueryOption) chromedp.QueryAction
type produceFileUserActionFunc = func(*[]byte) chromedp.Action
type usingVariableFunc = func(msg *json.RawMessage) chromedp.Action
type usingStringVariableFunc = func(txt *string) chromedp.Action

func bindTask(task config.Task) chromedp.Action {
	switch t := task.(type) {
	case *config.Navigate:
		return chromedp.ActionFunc(func(c context.Context) error {
			v, err := evalContext(c, t.URL)
			if err != nil {
				return nil
			}
			url := v.AsString()
			return tasks(chromedp.Navigate(url), printf("Navigate to `%s'", url)).Do(c)
		})
	case *config.NavigateForward:
		return tasks(chromedp.NavigateForward(), printf("Navigate forward"))
	case *config.NavigateBack:
		return tasks(chromedp.NavigateBack(), printf("Navigate back"))
	case *config.WaitVisible:
		return bindSelector("Wait until visible", chromedp.WaitVisible, t.Selector, t.Selectors, t.Options)
	case *config.Click:
		return bindSelector("Click", chromedp.Click, t.Selector, t.Selectors, t.Options)
	case *config.DoubleClick:
		return bindSelector("Double click", chromedp.DoubleClick, t.Selector, t.Selectors, t.Options)
	case *config.Blur:
		return bindSelector("Blur", chromedp.Blur, t.Selector, t.Selectors, t.Options)
	case *config.Clear:
		return bindSelector("Clear", chromedp.Clear, t.Selector, t.Selectors, t.Options)
	case *config.Sleep:
		return tasks(chromedp.Sleep(t.Duration), printf("Sleep %v", t.Duration))
	case *config.Reload:
		return tasks(chromedp.Reload(), printf("Reload"))
	case *config.Stop:
		return tasks(chromedp.Stop(), printf("Stop"))
	case *config.Screenshot:
		filename := cmp.Or(t.Name, "screenshot.png")
		if t.Selector == "" && len(t.Selectors) == 0 {
			return tasks(requestOutputFile(filename, chromedp.CaptureScreenshot), printf("Capture screenshot"))
		}

		return bindSelector("Capture screenshot", func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
			curry := func(file *[]byte) chromedp.Action {
				return chromedp.Screenshot(sel, file, opts...)
			}
			if t.Scale > 0 {
				curry = func(file *[]byte) chromedp.Action {
					return chromedp.ScreenshotScale(sel, t.Scale, file, opts...)
				}
			}

			return requestOutputFile(filename, curry)
		}, t.Selector, t.Selectors, t.Options)

	case *config.Eval:
		return usingVariable(t.Name, func(msg *json.RawMessage) chromedp.Action {
			return chromedp.ActionFunc(func(c context.Context) error {
				err := tasks(chromedp.Evaluate(t.Script, msg), printf("Evaluate script `%s'", t.Name)).Do(c)
				evalContextFrom(c).Variables[t.Name] = umarshalData(*msg)
				return err
			})
		})

	case *config.Title:
		return usingStringVariable(t.Name, func(str *string) chromedp.Action {
			return chromedp.ActionFunc(func(c context.Context) error {
				err := tasks(chromedp.Title(str), printf("Extract title into variable `%s'", t.Name)).Do(c)
				v, _ := gocty.ToCtyValue(*str, cty.String)
				evalContextFrom(c).Variables[t.Name] = v
				return err
			})
		})

	case *config.Version:
		return printBrowserVersion(browser.GetVersion())

	default:
		panic(fmt.Errorf("unexpected task type %T", t))
	}
}

func printBrowserVersion(params *browser.GetVersionParams) TaskFunc {
	return func(ctx context.Context) error {
		protocolVersion, product, revision, userAgent, jsVersion, err := params.Do(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Browser versions: %s/%s, %s, protocol=%s, js=%s\n",
			product, revision, userAgent, protocolVersion, jsVersion)
		return nil
	}
}

func umarshalData(msg json.RawMessage) cty.Value {
	// TODO Additional types might occur
	m := map[string]string{}
	_ = json.Unmarshal(msg, &m)
	v, _ := gocty.ToCtyValue(m, cty.Map(cty.String))
	return v
}

func bindSelector(desc string, fn produceQueryActionFunc, sel string, sels []*config.Selector, options *config.Options) Task {
	if sel != "" {
		sels = append(sels, &config.Selector{
			Target: sel,
			By:     config.BySearch,
		})
	}

	var tasks chromedp.Tasks = make([]chromedp.Action, len(sels))
	selectorDesc := make([]string, len(sels))
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
		selectorDesc[i] = fmt.Sprintf(
			"%s (by=%s, on=%s, at_least=%v, retry_interval=%v)", s.Target, s.By, s.On, options.AtLeast, options.RetryInterval)
	}
	tasks = append(tasks, printf("%s %s", desc, strings.Join(selectorDesc, ",")))
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

func bindAutomation(automation *config.Automation) []Task {
	actions := make([]Task, 0)
	for _, t := range automation.Tasks {
		actions = append(actions, bindTask(t))
	}
	return actions
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

func usingVariable(name string, fn usingVariableFunc) chromedp.Action {
	return chromedp.ActionFunc(func(c context.Context) error {
		res := mustAutomationResult(c)
		var msg json.RawMessage
		res.Outputs[name] = &msg
		return fn(&msg).Do(c)
	})
}

func usingStringVariable(name string, fn usingStringVariableFunc) chromedp.Action {
	return chromedp.ActionFunc(func(c context.Context) error {
		res := mustAutomationResult(c)
		var str string

		err := fn(&str).Do(c)

		var msg json.RawMessage
		msg, _ = json.Marshal(str)

		res.Outputs[name] = &msg

		return err
	})
}

func tasks(t ...Task) Tasks {
	return Tasks(t)
}

func printf(format string, args ...any) TaskFunc {
	return func(c context.Context) error {
		fmt.Printf(format, args...)
		fmt.Println()
		return nil
	}
}
