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

var _ = Describe("SelectorSet", func() {

	Describe("Set", func() {

		It("appends a selector with the text", func() {
			actual := new(internalcli.SelectorSet)
			err := actual.Set("button")

			Expect(err).NotTo(HaveOccurred())
			Expect(actual.Selectors).To(Equal([]model.Selector{
				{Target: "button", By: model.ByQueryAll},
			}))
		})

		It("splits on commas", func() {
			actual := new(internalcli.SelectorSet)
			err := actual.Set("button,input,a")

			Expect(err).NotTo(HaveOccurred())
			Expect(actual.Selectors).To(Equal([]model.Selector{
				{Target: "button", By: model.ByQueryAll},
				{Target: "input", By: model.ByQueryAll},
				{Target: "a", By: model.ByQueryAll},
			}))
		})

		It("defaults the strategy to QUERY_ALL", func() {
			actual := new(internalcli.SelectorSet)
			_ = actual.Set("button")

			Expect(actual.Selectors[0].By).To(Equal(model.ByQueryAll))
		})

		It("copies the strategy to each selector", func() {
			actual := &internalcli.SelectorSet{By: model.ByID}
			err := actual.Set("first,second")

			Expect(err).NotTo(HaveOccurred())
			Expect(actual.Selectors).To(Equal([]model.Selector{
				{Target: "first", By: model.ByID},
				{Target: "second", By: model.ByID},
			}))
		})

		It("appends across multiple calls", func() {
			actual := new(internalcli.SelectorSet)
			Expect(actual.Set("button")).To(Succeed())
			Expect(actual.Set("input")).To(Succeed())

			Expect(actual.Selectors).To(Equal([]model.Selector{
				{Target: "button", By: model.ByQueryAll},
				{Target: "input", By: model.ByQueryAll},
			}))
		})
	})

	Describe("String", func() {
		It("joins the selector targets with commas", func() {
			actual := new(internalcli.SelectorSet)
			_ = actual.Set("button,input")

			Expect(actual.String()).To(Equal("button,input"))
		})
	})

	Describe("integration with an App flag", func() {

		It("parses selectors from the command line", func() {
			set := new(internalcli.SelectorSet)
			app := &cli.App{
				Name: "app",
				Flags: []*cli.Flag{
					{
						Name:  "selector",
						Value: set,
					},
				},
				Action: func() {},
			}

			err := app.RunContext(context.Background(), []string{
				"app", "--selector", "button,input",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(set.Selectors).To(Equal([]model.Selector{
				{Target: "button", By: model.ByQueryAll},
				{Target: "input", By: model.ByQueryAll},
			}))
		})

		It("accumulates across repeated flag occurrences using the configured strategy", func() {
			set := &internalcli.SelectorSet{By: model.ByID}
			app := &cli.App{
				Name: "app",
				Flags: []*cli.Flag{
					{
						Name:  "selector",
						Value: set,
					},
				},
				Action: func() {},
			}

			err := app.RunContext(context.Background(), []string{
				"app", "--selector", "first", "--selector", "second,third",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(set.Selectors).To(Equal([]model.Selector{
				{Target: "first", By: model.ByID},
				{Target: "second", By: model.ByID},
				{Target: "third", By: model.ByID},
			}))
		})
	})
})
