// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	internalcli "github.com/Carbonfrost/autogun/pkg/internal/cli"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Protocol", func() {

	Describe("Set", func() {

		DescribeTable("examples",
			func(arg string, expected internalcli.Protocol) {
				actual := new(internalcli.Protocol)
				err := actual.Set(arg)

				Expect(err).NotTo(HaveOccurred())
				Expect(*actual).To(Equal(internalcli.Protocol(expected)))
			},
			Entry("chromedp", "chromedp", internalcli.Chromedp),
		)

		DescribeTable("errors",
			func(arg string, expected string) {
				actual := new(internalcli.Protocol)
				err := actual.Set(arg)

				Expect(err).To(MatchError(expected))
			},
			Entry("error", "unknown", "invalid value: \"unknown\""),
		)
	})

	Describe("String", func() {
		DescribeTable("examples",
			func(arg internalcli.Protocol, expected string) {
				Expect(internalcli.Protocol(arg).String()).To(Equal(expected))
			},
			Entry("chromedp", internalcli.Chromedp, "chromedp"),
		)
	})
})
