// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"

	"github.com/Carbonfrost/autogun/pkg/config"
)

func fromConfigFile(file *config.File) []*Automation {
	if file == nil {
		return nil
	}
	autos := make([]*Automation, 0, len(file.Automations))
	for _, a := range file.Automations {
		autos = append(autos, FromConfig(a))
	}
	return autos
}

// FromConfig converts a configuration automation into its model
// representation
func FromConfig(cfg *config.Automation) *Automation {
	if cfg == nil {
		return nil
	}
	tasks := make([]Task, 0, len(cfg.Tasks))
	for _, t := range cfg.Tasks {
		tasks = append(tasks, taskFromConfig(t))
	}
	return &Automation{
		Name:  cfg.Name,
		Tasks: tasks,
	}
}

func taskFromConfig(task config.Task) Task {
	switch t := task.(type) {
	case *config.Navigate:
		return &Navigate{
			Name: t.Name,
			URL:  ExpressionFromHCL(t.URL),
		}
	case *config.NavigateForward:
		return &NavigateForward{}
	case *config.NavigateBack:
		return &NavigateBack{}
	case *config.Title:
		return &Title{Name: t.Name}
	case *config.Eval:
		return &Eval{Name: t.Name, Script: t.Script}
	case *config.Blur:
		return &Blur{
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.Clear:
		return &Clear{
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.Click:
		return &Click{
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.DoubleClick:
		return &DoubleClick{
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.WaitVisible:
		return &WaitVisible{
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.Screenshot:
		return &Screenshot{
			Name:      t.Name,
			Scale:     t.Scale,
			Selectors: selectorsFromConfig(t.Selector, t.Selectors),
			Options:   optionsFromConfig(t.Options),
		}
	case *config.Sleep:
		return &Sleep{Duration: t.Duration}
	case *config.Reload:
		return &Reload{}
	case *config.Stop:
		return &Stop{}
	case *config.Version:
		return &Version{}
	default:
		panic(fmt.Errorf("unexpected task type %T", t))
	}
}

func selectorsFromConfig(selector string, sels []*config.Selector) []*Selector {
	out := make([]*Selector, 0, len(sels)+1)
	for _, s := range sels {
		out = append(out, selectorFromConfig(s))
	}
	if selector != "" {
		out = append(out, &Selector{
			Target: selector,
			By:     BySearch,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func selectorFromConfig(s *config.Selector) *Selector {
	if s == nil {
		return nil
	}
	return &Selector{
		Target: s.Target,
		By:     SelectorBy(s.By),
		On:     SelectorOn(s.On),
	}
}

func optionsFromConfig(o *config.Options) *Options {
	if o == nil {
		return nil
	}
	return &Options{
		RetryInterval: o.RetryInterval,
		AtLeast:       o.AtLeast,
	}
}
