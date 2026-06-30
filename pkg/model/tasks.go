// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import "time"

// Task is the basis of a step within an automation. Each task type mirrors the
// corresponding type in the config package, but without the HCL-specific
// declaration ranges and with HCL expressions replaced by the [Expression]
// indirection.
type Task interface {
	taskSigil()
}

type Navigate struct {
	Name string
	URL  Expression
}

type NavigateForward struct{}

type NavigateBack struct{}

type Title struct {
	Name string
}

type Eval struct {
	Name   string
	Script string
}

type Blur struct {
	Selectors []*Selector
	Options   *Options
}

type Clear struct {
	Selectors []*Selector
	Options   *Options
}

type Click struct {
	Selectors []*Selector
	Options   *Options
}

type DoubleClick struct {
	Selectors []*Selector
	Options   *Options
}

type WaitVisible struct {
	Selectors []*Selector
	Options   *Options
}

type Screenshot struct {
	Name      string
	Scale     float64
	Selectors []*Selector
	Options   *Options
}

type Options struct {
	RetryInterval *time.Duration
	AtLeast       *int
}

type Sleep struct {
	Duration time.Duration
}

type Reload struct{}

type Stop struct{}

type Version struct{}

func (*Blur) taskSigil()            {}
func (*Clear) taskSigil()           {}
func (*Click) taskSigil()           {}
func (*DoubleClick) taskSigil()     {}
func (*Eval) taskSigil()            {}
func (*Navigate) taskSigil()        {}
func (*NavigateBack) taskSigil()    {}
func (*NavigateForward) taskSigil() {}
func (*Reload) taskSigil()          {}
func (*Screenshot) taskSigil()      {}
func (*Sleep) taskSigil()           {}
func (*Stop) taskSigil()            {}
func (*Title) taskSigil()           {}
func (*Version) taskSigil()         {}
func (*WaitVisible) taskSigil()     {}
