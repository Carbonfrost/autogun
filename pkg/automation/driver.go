// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/chromedp/chromedp"
)

type Driver struct {
	automations map[string]*Automation
	allocator   *Allocator
	protocol    Protocol
	model       *model.Model
	logging     bool
}

// Automation get the automation by name
func (d *Driver) Automation(name string) *Automation {
	return d.automations[name]
}

// Execute will execute the named automation
func (d *Driver) Execute(ctx context.Context, auto *model.Automation) (*Result, error) {
	res := newResult()
	ctx, cancel, err := d.allocator.newContext(
		withAutomationResult(ctx, res),
	)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var emulate Task = TaskFunc(nil)
	if dev, ok := d.allocator.resolveDevice(); ok {
		fmt.Fprintf(os.Stderr, "Emulating device %s (%s)\n", dev.Name, d.allocator.DeviceID)
		emulate = chromedp.Emulate(bindDevice(dev))
	}

	a, err := d.buildAutomation(auto)
	if err != nil {
		return nil, err
	}

	return res, chromedp.Run(ctx, emulate, a)
}

func (d *Driver) buildAutomation(m *model.Automation) (*Automation, error) {
	actions := make([]Task, 0)
	for _, t := range m.Tasks {
		switch task := t.(type) {
		case *model.Flow:
			actions = append(actions, d.flow(task.Name))
		case *model.Source:
			actions = append(actions, d.runSource(task.Filename))
		default:
			tsk, err := d.protocol.BindTask(t)
			if err != nil {
				return nil, err
			}
			if d.logging {
				if desc := describe(t); desc != "" {
					tsk = tasks(Print(desc), tsk)
				}
			}
			actions = append(actions, tsk)
		}
	}

	return &Automation{
		Name:  m.Name,
		Tasks: actions,
	}, nil
}

func describe(task model.Task) string {
	switch t := task.(type) {
	case *model.Navigate:
		return fmt.Sprintf("Navigate to %q", t.URL)
	case *model.NavigateForward:
		return "Navigate forward"
	case *model.NavigateBack:
		return "Navigate back"
	case *model.WaitVisible:
		return describeSelector("Wait until visible", t.Selectors, t.Options)
	case *model.Click:
		return describeSelector("Click", t.Selectors, t.Options)
	case *model.DoubleClick:
		return describeSelector("Double click", t.Selectors, t.Options)
	case *model.Blur:
		return describeSelector("Blur", t.Selectors, t.Options)
	case *model.Clear:
		return describeSelector("Clear", t.Selectors, t.Options)
	case *model.SendKeys:
		return describeSelector("Send keys", t.Selectors, t.Options)
	case *model.SetValue:
		return describeSelector("Set value", t.Selectors, t.Options)
	case *model.Sleep:
		return fmt.Sprintf("Sleep %v", t.Duration)
	case *model.Reload:
		return "Reload"
	case *model.Stop:
		return "Stop"
	case *model.Screenshot:
		if len(t.Selectors) == 0 {
			return "Capture screenshot"
		}
		return describeSelector("Capture screenshot", t.Selectors, t.Options)
	case *model.Eval:
		return fmt.Sprintf("Evaluate script `%s'", t.Name)
	case *model.Title:
		return fmt.Sprintf("Extract title into variable `%s'", t.Name)
	case *model.InnerHTML:
		return describeSelector(fmt.Sprintf("Extract inner HTML into variable `%s'", t.Name), t.Selectors, t.Options)
	case *model.Version:
		return "Get browser version"
	default:
		return ""
	}
}

func describeSelector(desc string, sels []*model.Selector, options *model.Options) string {
	selectorDesc := make([]string, len(sels))
	var opts string
	if options != nil {
		opts = fmt.Sprintf(" at %v/%v", options.AtLeast, options.RetryInterval)
	}
	for i, s := range sels {
		selectorDesc[i] = fmt.Sprintf("%s (%s, on=%s%v)", s.Target, s.By, s.On, opts)
	}
	return fmt.Sprintf("%s %s", desc, strings.Join(selectorDesc, ","))
}

func (d *Driver) flow(name string) Task {
	return taskThunk(func(c context.Context) (Task, error) {
		a := d.Automation(name)
		if a == nil {
			return nil, fmt.Errorf("automation not found %q", name)
		}
		return a, nil
	})
}

func (d *Driver) runSource(source string) Task {
	return taskThunk(func(c context.Context) (Task, error) {
		root := os.DirFS(".")
		p := config.NewParser(root)
		file, diag := p.LoadFile(source)
		if diag.HasErrors() {
			return nil, diag
		}

		// TODO Assumes one automation per file which may not be true
		return d.buildAutomation(model.FromConfig(file.Automations[0]))
	})
}
