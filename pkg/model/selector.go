// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

type Selector struct {
	Target string
	By     SelectorBy
	On     SelectorOn
}

type SelectorBy string
type SelectorOn string

const (
	BySearch   SelectorBy = "SEARCH"
	ByJSPath   SelectorBy = "JS_PATH"
	ByID       SelectorBy = "ID"
	ByQuery    SelectorBy = "QUERY"
	ByQueryAll SelectorBy = "QUERY_ALL"
)

const (
	OnReady      SelectorOn = "READY"
	OnVisible    SelectorOn = "VISIBLE"
	OnNotVisible SelectorOn = "NOT_VISIBLE"
	OnEnabled    SelectorOn = "ENABLED"
	OnSelected   SelectorOn = "SELECTED"
	OnNotPresent SelectorOn = "NOT_PRESENT"
)
