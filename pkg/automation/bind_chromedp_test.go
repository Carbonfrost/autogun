// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation_test

import (
	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bind", func() {

	Describe("bind Task", func() {

		DescribeTable("examples",
			func(task model.Task) {
				var (
					driver *automation.Driver
					output *automation.Automation
					err    error
				)
				bind := func() {
					driver, err = automation.Bind(&model.Model{
						Automations: []*model.Automation{
							{
								Name: "a",
								Tasks: []model.Task{
									task,
								},
							},
						},
					})
					output = driver.Automation("a")
				}

				Expect(bind).NotTo(Panic())
				Expect(err).NotTo(HaveOccurred())
				Expect(output.Tasks).To(HaveLen(1))
				Expect(output.Tasks[0]).NotTo(BeNil())
			},
			Entry("click", new(model.Click)),
			Entry("double_click", new(model.DoubleClick)),
			Entry("blur", new(model.Blur)),
			Entry("clear", new(model.Clear)),
			Entry("eval", new(model.Eval)),
			Entry("navigate", new(model.Navigate)),
			Entry("navigate_back", new(model.NavigateBack)),
			Entry("navigate_forward", new(model.NavigateForward)),
			Entry("reload", new(model.Reload)),
			Entry("screenshot", new(model.Screenshot)),
			Entry("send_keys", new(model.SendKeys)),
			Entry("sleep", new(model.Sleep)),
			Entry("stop", new(model.Stop)),
			Entry("title", new(model.Title)),
			Entry("wait_visible", new(model.WaitVisible)),
		)
	})
})
