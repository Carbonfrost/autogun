// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace_test

import (
	"context"

	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/Carbonfrost/autogun/pkg/workspace"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/expr"
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
		Entry(nil, "blur"),
		Entry(nil, "clear"),
		Entry(nil, "click"),
		Entry(nil, "doubleclick"),
		Entry(nil, "eval"),
		Entry(nil, "inner_html"),
		Entry(nil, "navigate"),
		Entry(nil, "back"),
		Entry(nil, "forward"),
		Entry(nil, "options"),
		Entry(nil, "reload"),
		Entry(nil, "screenshot"),
		Entry(nil, "select"),
		Entry(nil, "send_keys"),
		Entry(nil, "sleep"),
		Entry(nil, "stop"),
		Entry(nil, "title"),
		Entry(nil, "version"),
		Entry(nil, "wait_visible"),
	)
})

var _ = Describe("select propagation", func() {

	var evaluate = func(args ...string) []model.Task {
		query := &workspace.AutomationQuery{Automation: &model.Automation{}}
		app := &cli.App{
			Name: "app",
			Args: []*cli.Arg{
				{
					Name:  "sources",
					Value: cli.List(),
					NArg:  cli.TakeUntilNextFlag,
				},
				{
					Name: "expression",
					Value: &expr.Expression{
						Exprs: workspace.Exprs(),
					},
					Options: cli.SortedExprs,
				},
			},
			Action: func(c *cli.Context) error {
				e := c.Value("expression").(*expr.Expression)
				return e.Evaluate(c, query)
			},
		}
		err := app.RunContext(context.Background(), append([]string{"app", "."}, args...))
		Expect(err).NotTo(HaveOccurred())
		return query.Automation.Tasks
	}

	It("propagates the selector set to a click task", func() {
		tasks := evaluate("-select", "button", "-click")
		Expect(tasks).To(HaveLen(1))
		Expect(tasks[0]).To(Equal(&model.Click{
			Selectors: []*model.Selector{
				{Target: "button", By: model.ByQueryAll},
			},
		}))
	})

	It("propagates the selector set to each subsequent task", func() {
		tasks := evaluate("-select", "input", "-clear", "-click", "-blur")
		Expect(tasks).To(HaveLen(3))
		want := []*model.Selector{{Target: "input", By: model.ByQueryAll}}
		Expect(tasks[0]).To(Equal(&model.Clear{Selectors: want}))
		Expect(tasks[1]).To(Equal(&model.Click{Selectors: want}))
		Expect(tasks[2]).To(Equal(&model.Blur{Selectors: want}))
	})

	It("splits the selector set on commas", func() {
		tasks := evaluate("-select", "a,b", "-wait_visible")
		Expect(tasks).To(HaveLen(1))
		Expect(tasks[0]).To(Equal(&model.WaitVisible{
			Selectors: []*model.Selector{
				{Target: "a", By: model.ByQueryAll},
				{Target: "b", By: model.ByQueryAll},
			},
		}))
	})

	It("replaces the selector set when -select is used again", func() {
		tasks := evaluate("-select", "first", "-click", "-select", "second", "-click")
		Expect(tasks).To(HaveLen(2))
		Expect(tasks[0]).To(Equal(&model.Click{
			Selectors: []*model.Selector{{Target: "first", By: model.ByQueryAll}},
		}))
		Expect(tasks[1]).To(Equal(&model.Click{
			Selectors: []*model.Selector{{Target: "second", By: model.ByQueryAll}},
		}))
	})

	It("lets the local screenshot selector take precedence over the set", func() {
		tasks := evaluate("-select", "fromset", "-screenshot", "selector=local")
		Expect(tasks).To(HaveLen(1))
		shot := tasks[0].(*model.Screenshot)
		Expect(shot.Selectors).To(Equal([]*model.Selector{
			{Target: "local"},
		}))
	})

	It("falls back to the selector set for a screenshot without a local selector", func() {
		tasks := evaluate("-select", "fromset", "-screenshot")
		Expect(tasks).To(HaveLen(1))
		shot := tasks[0].(*model.Screenshot)
		Expect(shot.Selectors).To(Equal([]*model.Selector{
			{Target: "fromset", By: model.ByQueryAll},
		}))
	})
})
