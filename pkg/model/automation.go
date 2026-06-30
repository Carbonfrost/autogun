// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Automation is a multi-step automated process.
type Automation struct {
	Name  string
	Tasks []Task
}
