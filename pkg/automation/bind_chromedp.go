// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"

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
		return taskThunk(func(c context.Context) (Task, error) {
			v, err := evalContext(c, t.URL)
			if err != nil {
				return nil, err
			}
			return chromedp.Navigate(v.AsString()), nil
		})
	case *model.NavigateForward:
		return chromedp.NavigateForward()
	case *model.NavigateBack:
		return chromedp.NavigateBack()
	case *model.WaitVisible:
		return bindSelector(chromedp.WaitVisible, t.Selectors, t.Options)
	case *model.Click:
		return bindSelector(chromedp.Click, t.Selectors, t.Options)
	case *model.DoubleClick:
		return bindSelector(chromedp.DoubleClick, t.Selectors, t.Options)
	case *model.Blur:
		return bindSelector(chromedp.Blur, t.Selectors, t.Options)
	case *model.Clear:
		return bindSelector(chromedp.Clear, t.Selectors, t.Options)
	case *model.SendKeys:
		return taskThunk(func(c context.Context) (Task, error) {
			v, err := evalContext(c, t.Keys)
			if err != nil {
				return nil, err
			}
			keys := v.AsString()
			return bindSelector(func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
				return chromedp.SendKeys(sel, keys, opts...)
			}, t.Selectors, t.Options), nil
		})
	case *model.SetValue:
		return taskThunk(func(c context.Context) (Task, error) {
			v, err := evalContext(c, t.Value)
			if err != nil {
				return nil, err
			}
			value := v.AsString()
			return bindSelector(func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
				return chromedp.SetValue(sel, value, opts...)
			}, t.Selectors, t.Options), nil
		})
	case *model.Sleep:
		return chromedp.Sleep(t.Duration)
	case *model.Reload:
		return chromedp.Reload()
	case *model.Stop:
		return chromedp.Stop()
	case *model.Screenshot:
		filename := cmp.Or(t.Name, "screenshot.png")
		if len(t.Selectors) == 0 {
			return requestOutputFile(filename, chromedp.CaptureScreenshot)
		}

		return bindSelector(func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
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
			return chromedp.Evaluate(t.Script, msg)
		})

	case *model.Title:
		return usingStringVariable(t.Name, func(str *string) chromedp.Action {
			return chromedp.Title(str)
		})

	case *model.InnerHTML:
		return usingStringVariable(t.Name, func(str *string) chromedp.Action {
			return bindSelector(func(sel any, opts ...chromedp.QueryOption) chromedp.QueryAction {
				return chromedp.InnerHTML(sel, str, opts...)
			}, t.Selectors, t.Options)
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

func bindSelector(fn produceQueryActionFunc, sels []*model.Selector, options *model.Options) Task {
	var tasks chromedp.Tasks = make([]chromedp.Action, len(sels))
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
		err := fn(&msg).Do(c)

		evalContextFrom(c).Variables[name] = umarshalData(msg)
		return err
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

		v, _ := gocty.ToCtyValue(&str, cty.String)
		evalContextFrom(c).Variables[name] = v

		return err
	})
}

func tasks(t ...Task) Tasks {
	return Tasks(t)
}
