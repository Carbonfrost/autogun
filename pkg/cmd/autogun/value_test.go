// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun_test

import (
	"github.com/Carbonfrost/autogun/pkg/cmd/autogun"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Engine", func() {

	Describe("Set", func() {

		DescribeTable("examples",
			func(arg string, expected autogun.Engine) {
				actual := new(autogun.Engine)
				err := actual.Set(arg)

				Expect(err).NotTo(HaveOccurred())
				Expect(*actual).To(Equal(autogun.Engine(expected)))
			},
			Entry("chromedp", "chromedp", autogun.Chromedp),
		)

		DescribeTable("errors",
			func(arg string, expected string) {
				actual := new(autogun.Engine)
				err := actual.Set(arg)

				Expect(err).To(MatchError(expected))
			},
			Entry("error", "unknown", "invalid value: \"unknown\""),
		)
	})

	Describe("String", func() {
		DescribeTable("examples",
			func(arg autogun.Engine, expected string) {
				Expect(autogun.Engine(arg).String()).To(Equal(expected))
			},
			Entry("chromedp", autogun.Chromedp, "chromedp"),
		)
	})
})
