// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model_test

import (
	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/Carbonfrost/autogun/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromConfig", func() {

	Describe("task conversion", func() {

		DescribeTable("converts each config task type",
			func(in config.Task, expected model.Task) {
				out := model.FromConfig(&config.Automation{
					Tasks: []config.Task{in},
				})

				Expect(out.Tasks).To(HaveLen(1))
				Expect(out.Tasks[0]).To(BeAssignableToTypeOf(expected))
			},
			Entry("blur", new(config.Blur), new(model.Blur)),
			Entry("clear", new(config.Clear), new(model.Clear)),
			Entry("click", new(config.Click), new(model.Click)),
			Entry("double_click", new(config.DoubleClick), new(model.DoubleClick)),
			Entry("eval", new(config.Eval), new(model.Eval)),
			Entry("navigate", new(config.Navigate), new(model.Navigate)),
			Entry("navigate_back", new(config.NavigateBack), new(model.NavigateBack)),
			Entry("navigate_forward", new(config.NavigateForward), new(model.NavigateForward)),
			Entry("reload", new(config.Reload), new(model.Reload)),
			Entry("screenshot", new(config.Screenshot), new(model.Screenshot)),
			Entry("sleep", new(config.Sleep), new(model.Sleep)),
			Entry("stop", new(config.Stop), new(model.Stop)),
			Entry("title", new(config.Title), new(model.Title)),
			Entry("version", new(config.Version), new(model.Version)),
			Entry("wait_visible", new(config.WaitVisible), new(model.WaitVisible)),
		)
	})
})
