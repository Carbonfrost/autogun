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
	apply(*Driver)
}

type Protocol interface { // Engine
	BindTask(model.Task) (Task, error)
	NewExecAllocator(parent context.Context, opts *AllocatorOptions) (context.Context, context.CancelFunc, error)
	NewRemoteAllocator(parent context.Context, url string, opts *AllocatorOptions) (context.Context, context.CancelFunc, error)
}

// SupportedProtocol is one of the supported binders
type SupportedProtocol int

const (
	// ProtocolChromedp is a Binder for using Chrome DevTools Protocol and headless
	// Chrome/Chromium to run the automation
	ProtocolChromedp SupportedProtocol = iota
)

var errNotSupportedProtocol = errors.New("unsupported binder")

// Bind converts the model into an automation driver
func Bind(m *model.Model, opts ...Option) (*Driver, error) {
	b := &Driver{
		model:     m,
		protocol:  ProtocolChromedp,
		allocator: &Allocator{},
	}
	for _, o := range opts {
		o.apply(b)
	}
	automations := map[string]*Automation{}

	var err error

	// TODO Don't eagerly build these automations; should be handled
	// within Driver.flow and possibly cache
	for _, auto := range m.Automations {
		automations[auto.Name], err = b.buildAutomation(auto)
		if err != nil {
			return nil, err
		}
	}
	b.automations = automations
	return b, nil
}

func (s SupportedProtocol) BindTask(cfg model.Task) (Task, error) {
	switch s {
	case ProtocolChromedp:
		return bindTask(cfg), nil
	default:
		return nil, errNotSupportedProtocol
	}
}

func (s SupportedProtocol) NewExecAllocator(parent context.Context, opts *AllocatorOptions) (context.Context, context.CancelFunc, error) {
	switch s {
	case ProtocolChromedp:
		if opts == nil {
			opts = &AllocatorOptions{}
		}
		opts.warnRemoteOnlyForExec()

		// Preserve the behavior of chromedp.NewContext, which applies
		// DefaultExecAllocatorOptions when the parent has no allocator,
		// then layer the caller's decomposed options on top.
		execOpts := append(
			append([]chromedp.ExecAllocatorOption{}, chromedp.DefaultExecAllocatorOptions[:]...),
			opts.execOptions()...,
		)
		allocatorContext, cancelAllocator := chromedp.NewExecAllocator(parent, execOpts...)
		res, cancelInner := chromedp.NewContext(allocatorContext)
		return res, func() {
			defer cancelAllocator()
			defer cancelInner()
		}, nil
	default:
		return nil, nil, errNotSupportedProtocol
	}
}

func (s SupportedProtocol) NewRemoteAllocator(parent context.Context, url string, opts *AllocatorOptions) (context.Context, context.CancelFunc, error) {
	switch s {
	case ProtocolChromedp:
		if opts == nil {
			opts = &AllocatorOptions{}
		}
		opts.warnExecOnlyForRemote()

		allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(parent, url, opts.remoteOptions()...)
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

func WithAllocator(v *Allocator) Option {
	return optionFunc(func(a *Driver) {
		a.allocator = v
	})
}

func WithProtocol(p Protocol) Option {
	return optionFunc(func(a *Driver) {
		a.protocol = p
	})
}

type optionFunc func(*Driver)

func (f optionFunc) apply(a *Driver) {
	f(a)
}
