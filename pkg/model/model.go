// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
)

// Model is the collection of automations available in the workspace, bound and
// ready to execute.
type Model struct {
	Automations []*Automation
}

// New creates a model from the given configuration files
func New(files ...*config.File) *Model {
	m := &Model{}
	for _, file := range files {
		for _, auto := range fromConfigFile(file) {
			m.Automations = append(m.Automations, auto)
		}
	}
	return m
}

// Automation retrieves the automation by name
func (m *Model) Automation(name string) *Automation {
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
