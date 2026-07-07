// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace_test

import (
	"github.com/Carbonfrost/autogun/pkg/workspace"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exprs", func() {

	var names = func() []string {
		res := []string{}
		for _, e := range workspace.Exprs() {
			res = append(res, e.Name)
		}
		return res
	}()

	DescribeTable("contains expected operand", func(name string) {
		Expect(names).To(ContainElement(name))
	},
		EntryDescription("%[1]s"),
		XEntry(nil, "blur"),
		XEntry(nil, "clear"),
		XEntry(nil, "click"),
		XEntry(nil, "doubleclick"),
		Entry(nil, "eval"),
		Entry(nil, "navigate"),
		Entry(nil, "back"),
		Entry(nil, "forward"),
		XEntry(nil, "options"),
		Entry(nil, "reload"),
		Entry(nil, "screenshot"),
		Entry(nil, "sleep"),
		Entry(nil, "stop"),
		Entry(nil, "title"),
		Entry(nil, "version"),
		XEntry(nil, "wait_visible"),
	)
})
