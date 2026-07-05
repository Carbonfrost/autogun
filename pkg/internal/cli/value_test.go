// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	internalcli "github.com/Carbonfrost/autogun/pkg/internal/cli"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Engine", func() {

	Describe("Set", func() {

		DescribeTable("examples",
			func(arg string, expected internalcli.Engine) {
				actual := new(internalcli.Engine)
				err := actual.Set(arg)

				Expect(err).NotTo(HaveOccurred())
				Expect(*actual).To(Equal(internalcli.Engine(expected)))
			},
			Entry("chromedp", "chromedp", internalcli.Chromedp),
		)

		DescribeTable("errors",
			func(arg string, expected string) {
				actual := new(internalcli.Engine)
				err := actual.Set(arg)

				Expect(err).To(MatchError(expected))
			},
			Entry("error", "unknown", "invalid value: \"unknown\""),
		)
	})

	Describe("String", func() {
		DescribeTable("examples",
			func(arg internalcli.Engine, expected string) {
				Expect(internalcli.Engine(arg).String()).To(Equal(expected))
			},
			Entry("chromedp", internalcli.Chromedp, "chromedp"),
		)
	})
})
