// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"flag"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/model"
)

// SelectorSet is a flag.Value proxy for a slice of model.Selector.
type SelectorSet struct {
	// By sets up the selector strategy that applies to each selector. It can't
	// actually be set via the flag Value.
	By        model.SelectorBy
	Selectors []*model.Selector
}

func (s *SelectorSet) by() model.SelectorBy {
	if s.By == "" {
		return model.ByQueryAll
	}
	return s.By
}

func (s *SelectorSet) Value() []*model.Selector {
	return s.Selectors
}

// Set splits its argument on commas and appends a Selector per resulting text,
// copying the By strategy to each. When By is unset, it defaults to
// model.ByQueryAll.
func (s *SelectorSet) Set(arg string) error {
	for _, text := range strings.Split(arg, ",") {
		s.Selectors = append(s.Selectors, &model.Selector{
			Target: text,
			By:     s.by(),
		})
	}
	return nil
}

func (s *SelectorSet) Reset() {
	s.By = ""
	s.Selectors = nil
}

func (s *SelectorSet) String() string {
	texts := make([]string, len(s.Selectors))
	for i, sel := range s.Selectors {
		texts[i] = sel.Target
	}
	return strings.Join(texts, ",")
}

var _ flag.Value = (*SelectorSet)(nil)
