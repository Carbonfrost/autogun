// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package automation

import (
	"context"
)

// Task provides the basis of a step in the automation
type Task interface {
	// Do executes the action using the provided context and frame handler.
	Do(context.Context) error
}

// TaskFunc implements a task from function
type TaskFunc func(context.Context) error

// Automation is a multi-step automated process
type Automation struct {
	// Name gets the name of automation
	Name string

	// Tasks provides the tasks in the automation
	Tasks []Task
}

func (f TaskFunc) Do(c context.Context) error {
	if f == nil {
		return nil
	}
	return f(c)
}

func (a *Automation) Do(c context.Context) error {
	for _, task := range a.Tasks {
		err := task.Do(c)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ Task = (*Automation)(nil)
