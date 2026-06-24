// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"os"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"
)

// Model is the automations available
type Model struct {
	Automations []*automation.Automation
}

// New creates a Model by binding the automations declared in the given
// configuration files.
func New(files ...*config.File) *Model {
	m := &Model{}
	for _, file := range files {
		for _, auto := range file.Automations {
			a, err := automation.Bind(auto, automation.UsingChromedp)
			if err != nil {
				logError(err)
				continue
			}
			m.Automations = append(m.Automations, a)
		}
	}
	return m
}

// Automation retrieves the automation by name
func (m *Model) Automation(name string) *automation.Automation {
	for _, auto := range m.Automations {
		if auto.Name == name {
			return auto
		}
	}
	return nil
}

func logError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
