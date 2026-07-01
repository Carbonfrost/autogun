// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"errors"

	"github.com/Carbonfrost/autogun/pkg/model"
)

// Option provides an option to the automation builder
type Option func(*Automation)

type Protocol interface { // Engine
	BindAutomation(*model.Automation) (*Automation, error)
	BindTask(model.Task) (Task, error)
}

// SupportedProtocol is one of the supported binders
type SupportedProtocol int

const (
	// UsingChromedp is a Binder for using Chrome DevTools Protocol and headless
	// Chrome/Chromium to run the automation
	UsingChromedp SupportedProtocol = iota
)

var errNotSupportedProtocol = errors.New("unsupported binder")

// Bind converts the model into an automation
func Bind(m *model.Automation, using Protocol, opts ...Option) (*Automation, error) {
	if using == nil {
		using = UsingChromedp
	}
	auto, err := using.BindAutomation(m)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o(auto)
	}
	return auto, nil
}

func (s SupportedProtocol) BindAutomation(m *model.Automation) (*Automation, error) {
	switch s {
	case UsingChromedp:
		return &Automation{
			Name:  m.Name,
			Tasks: bindAutomation(m),
		}, nil
	default:
		return nil, errNotSupportedProtocol
	}
}

func (s SupportedProtocol) BindTask(cfg model.Task) (Task, error) {
	switch s {
	case UsingChromedp:
		return bindTask(cfg), nil
	default:
		return nil, errNotSupportedProtocol
	}
}
