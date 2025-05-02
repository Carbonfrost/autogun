// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package automation

import (
	"errors"

	"github.com/Carbonfrost/autogun/pkg/config"
)

// Option provides an option to the automation builder
type Option func(*Automation)

type Binder interface { // Engine
	BindAutomation(*config.Automation) (*Automation, error)
	BindTask(config.Task) (Task, error)
}

// SupportedBinder is one of the supported binders
type SupportedBinder int

const (
	// UsingChromedp is a Binder for using Chrome DevTools Protocol and headless
	// Chrome/Chromium to run the automation
	UsingChromedp SupportedBinder = iota
)

var errNotSupportedBinder = errors.New("unsupported binder")

// Bind converts the configuration into an automation
func Bind(cfg *config.Automation, using Binder, opts ...Option) (*Automation, error) {
	if using == nil {
		using = UsingChromedp
	}
	auto, err := using.BindAutomation(cfg)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o(auto)
	}
	return auto, nil
}

func (s SupportedBinder) BindAutomation(cfg *config.Automation) (*Automation, error) {
	switch s {
	case UsingChromedp:
		return &Automation{
			Name:  cfg.Name,
			Tasks: bindAutomation(cfg),
		}, nil
	default:
		return nil, errNotSupportedBinder
	}
}

func (s SupportedBinder) BindTask(cfg config.Task) (Task, error) {
	switch s {
	case UsingChromedp:
		return bindTask(cfg), nil
	default:
		return nil, errNotSupportedBinder
	}
}
