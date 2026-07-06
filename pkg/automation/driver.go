// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/chromedp/chromedp"
)

type Driver struct {
	automations map[string]*Automation
	allocator   *Allocator
	protocol    Protocol
	model       *model.Model
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
			actions = append(actions, tsk)
		}
	}

	return &Automation{
		Name:  m.Name,
		Tasks: actions,
	}, nil
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
