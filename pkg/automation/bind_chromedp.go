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

	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type produceQueryActionFunc = func(any, ...chromedp.QueryOption) chromedp.QueryAction
type produceFileUserActionFunc = func(*[]byte) chromedp.Action
type usingVariableFunc = func(msg *json.RawMessage) chromedp.Action
type usingStringVariableFunc = func(txt *string) chromedp.Action

func bindTask(task model.Task) chromedp.Action {
	switch t := task.(type) {
	case *model.Navigate:
		return chromedp.ActionFunc(func(c context.Context) error {
			v, err := evalContext(c, t.URL)
			if err != nil {
				return nil
			}
			url := v.AsString()
			return tasks(chromedp.Navigate(url), printf("Navigate to `%s'", url)).Do(c)
		})
	case *model.NavigateForward:
		return tasks(chromedp.NavigateForward(), printf("Navigate forward"))
	case *model.NavigateBack:
		return tasks(chromedp.NavigateBack(), printf("Navigate back"))
	case *model.WaitVisible:
		return bindSelector("Wait until visible", chromedp.WaitVisible, t.Selectors, t.Options)
	case *model.Click:
		return bindSelector("Click", chromedp.Click, t.Selectors, t.Options)
	case *model.DoubleClick:
		return bindSelector("Double click", chromedp.DoubleClick, t.Selectors, t.Options)
	case *model.Blur:
		return bindSelector("Blur", chromedp.Blur, t.Selectors, t.Options)
	case *model.Clear:
		return bindSelector("Clear", chromedp.Clear, t.Selectors, t.Options)
	case *model.Sleep:
		return tasks(chromedp.Sleep(t.Duration), printf("Sleep %v", t.Duration))
	case *model.Reload:
		return tasks(chromedp.Reload(), printf("Reload"))
	case *model.Stop:
		return tasks(chromedp.Stop(), printf("Stop"))
	case *model.Screenshot:
		filename := cmp.Or(t.Name, "screenshot.png")
		if len(t.Selectors) == 0 {
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
		}, t.Selectors, t.Options)

	case *model.Eval:
		return usingVariable(t.Name, func(msg *json.RawMessage) chromedp.Action {
			return chromedp.ActionFunc(func(c context.Context) error {
				err := tasks(chromedp.Evaluate(t.Script, msg), printf("Evaluate script `%s'", t.Name)).Do(c)
				evalContextFrom(c).Variables[t.Name] = umarshalData(*msg)
				return err
			})
		})

	case *model.Title:
		return usingStringVariable(t.Name, func(str *string) chromedp.Action {
			return chromedp.ActionFunc(func(c context.Context) error {
				err := tasks(chromedp.Title(str), printf("Extract title into variable `%s'", t.Name)).Do(c)
				v, _ := gocty.ToCtyValue(*str, cty.String)
				evalContextFrom(c).Variables[t.Name] = v
				return err
			})
		})

	case *model.Version:
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

func bindSelector(desc string, fn produceQueryActionFunc, sels []*model.Selector, options *model.Options) Task {
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

func bindSelectorBy(s model.SelectorBy) chromedp.QueryOption {
	switch s {
	case model.BySearch:
		return chromedp.BySearch
	case model.ByJSPath:
		return chromedp.ByJSPath
	case model.ByID:
		return chromedp.ByID
	case model.ByQuery:
		return chromedp.ByQuery
	case model.ByQueryAll:
		return chromedp.ByQueryAll
	}
	return nil
}

func bindSelectorOn(s model.SelectorOn) chromedp.QueryOption {
	switch s {
	case model.OnReady:
		return chromedp.NodeReady
	case model.OnVisible:
		return chromedp.NodeVisible
	case model.OnNotVisible:
		return chromedp.NodeNotVisible
	case model.OnEnabled:
		return chromedp.NodeEnabled
	case model.OnSelected:
		return chromedp.NodeSelected
	case model.OnNotPresent:
		return chromedp.NodeNotPresent
	}
	return nil
}

func bindQueryOptions(opts *model.Options) (results []chromedp.QueryOption) {
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

func bindAutomation(automation *model.Automation) []Task {
	actions := make([]Task, 0)
	for _, t := range automation.Tasks {
		actions = append(actions, bindTask(t))
	}
	return actions
}

func bindDevice(dev model.Device) chromedp.Device {
	return deviceImpl(dev)
}

type deviceImpl model.Device

func (d deviceImpl) Device() device.Info {
	return device.Info{
		Name:      d.Name,
		UserAgent: d.UserAgent,
		Width:     d.Width,
		Height:    d.Height,
		Scale:     d.Scale,
		Landscape: d.Landscape,
		Mobile:    d.Mobile,
		Touch:     d.Touch,
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
