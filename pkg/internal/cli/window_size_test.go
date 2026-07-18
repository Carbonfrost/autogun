// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"context"

	internalcli "github.com/Carbonfrost/autogun/pkg/internal/cli"
	"github.com/Carbonfrost/autogun/pkg/model"
	cli "github.com/Carbonfrost/joe-cli"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WindowSize", func() {

	Describe("Set", func() {

		DescribeTable("examples",
			func(arg string, expected model.WindowSize) {
				actual := new(internalcli.WindowSize)
				err := actual.Set(arg)

				Expect(err).NotTo(HaveOccurred())
				Expect(actual.WindowSize).To(Equal(expected))
			},
			Entry("WxH", "1280x720", model.WindowSize{Width: 1280, Height: 720}),
			Entry("W,H", "1280,720", model.WindowSize{Width: 1280, Height: 720}),
			Entry("trims surrounding spaces", " 1280 x 720 ", model.WindowSize{Width: 1280, Height: 720}),
			Entry("empty leaves the value unset", "", model.WindowSize{}),
		)

		DescribeTable("errors",
			func(arg string, expected string) {
				actual := new(internalcli.WindowSize)
				err := actual.Set(arg)

				Expect(err).To(MatchError(expected))
			},
			Entry("missing separator", "1280", `invalid window size "1280": expected WxH`),
			Entry("non-numeric width", "axb", `invalid window size "axb": strconv.Atoi: parsing "a": invalid syntax`),
			Entry("non-numeric height", "1280xb", `invalid window size "1280xb": strconv.Atoi: parsing "b": invalid syntax`),
		)

		It("overwrites on a subsequent call", func() {
			actual := new(internalcli.WindowSize)
			Expect(actual.Set("800x600")).To(Succeed())
			Expect(actual.Set("1024x768")).To(Succeed())

			Expect(actual.WindowSize).To(Equal(model.WindowSize{Width: 1024, Height: 768}))
		})
	})

	Describe("Get", func() {
		It("returns the parsed model.WindowSize", func() {
			actual := new(internalcli.WindowSize)
			Expect(actual.Set("1280x720")).To(Succeed())

			Expect(actual.Get()).To(Equal(model.WindowSize{Width: 1280, Height: 720}))
		})
	})

	Describe("String", func() {
		DescribeTable("examples",
			func(arg model.WindowSize, expected string) {
				actual := &internalcli.WindowSize{WindowSize: arg}
				Expect(actual.String()).To(Equal(expected))
			},
			Entry("WxH", model.WindowSize{Width: 1280, Height: 720}, "1280x720"),
			Entry("zero value", model.WindowSize{}, "0x0"),
		)
	})

	Describe("integration with an App flag", func() {

		It("parses the window size from the command line", func() {
			set := new(internalcli.WindowSize)
			app := &cli.App{
				Name: "app",
				Flags: []*cli.Flag{
					{
						Name:  "window-size",
						Value: set,
					},
				},
				Action: func() {},
			}

			err := app.RunContext(context.Background(), []string{
				"app", "--window-size", "1280x720",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(set.WindowSize).To(Equal(model.WindowSize{Width: 1280, Height: 720}))
		})

		It("reports parse errors from the command line", func() {
			app := &cli.App{
				Name: "app",
				Flags: []*cli.Flag{
					{
						Name:  "window-size",
						Value: new(internalcli.WindowSize),
					},
				},
				Action: func() {},
			}

			err := app.RunContext(context.Background(), []string{
				"app", "--window-size", "nonsense",
			})

			Expect(err).To(HaveOccurred())
		})
	})
})
