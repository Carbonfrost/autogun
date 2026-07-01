// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"errors"

	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/chromedp/chromedp"
)

// Option provides an option to the automation builder
type Option interface {
	apply(*automationBuilder)
}

type Protocol interface { // Engine
	BindAutomation(*model.Automation) (*Automation, error)
	BindTask(model.Task) (Task, error)
	NewExecAllocator(parent context.Context) (context.Context, context.CancelFunc, error)
	NewRemoteAllocator(parent context.Context, url string) (context.Context, context.CancelFunc, error)
}

// SupportedProtocol is one of the supported binders
type SupportedProtocol int

const (
	// ProtocolChromedp is a Binder for using Chrome DevTools Protocol and headless
	// Chrome/Chromium to run the automation
	ProtocolChromedp SupportedProtocol = iota
)

var errNotSupportedProtocol = errors.New("unsupported binder")

// Bind converts the model into an automation
func Bind(m *model.Automation, opts ...Option) (*Automation, error) {
	b := &automationBuilder{
		protocol: ProtocolChromedp,
	}
	for _, o := range opts {
		o.apply(b)
	}
	return b.build(m)
}

func (s SupportedProtocol) BindAutomation(m *model.Automation) (*Automation, error) {
	switch s {
	case ProtocolChromedp:
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
	case ProtocolChromedp:
		return bindTask(cfg), nil
	default:
		return nil, errNotSupportedProtocol
	}
}

func (s SupportedProtocol) NewExecAllocator(parent context.Context) (context.Context, context.CancelFunc, error) {
	switch s {
	case ProtocolChromedp:
		res, cancel := chromedp.NewContext(parent)
		return res, cancel, nil
	default:
		return nil, nil, errNotSupportedProtocol
	}
}

func (s SupportedProtocol) NewRemoteAllocator(parent context.Context, url string) (context.Context, context.CancelFunc, error) {
	switch s {
	case ProtocolChromedp:
		allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(parent, url)
		res, cancelInner := chromedp.NewContext(
			allocatorContext,
		)
		return res, func() {
			defer cancelAllocator()
			defer cancelInner()
		}, nil
	default:
		return nil, nil, errNotSupportedProtocol
	}
}

func WithProtocol(p Protocol) Option {
	return optionFunc(func(a *automationBuilder) {
		a.protocol = p
	})
}

type automationBuilder struct {
	protocol Protocol
}

func (a *automationBuilder) build(m *model.Automation) (*Automation, error) {
	auto, err := a.protocol.BindAutomation(m)
	if err != nil {
		return nil, err
	}
	return auto, nil
}

type optionFunc func(*automationBuilder)

func (f optionFunc) apply(a *automationBuilder) {
	f(a)
}
