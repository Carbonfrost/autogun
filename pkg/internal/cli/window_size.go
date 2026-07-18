// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/model"
)

// WindowSize is a flag.Value/flag.Getter proxy for model.WindowSize, parsed
// from "WxH" or "W,H".
type WindowSize struct {
	model.WindowSize
}

// Get implements flag.Getter, returning the parsed model.WindowSize.
func (s *WindowSize) Get() any {
	return s.WindowSize
}

// Set parses arg as "WxH" or "W,H" and stores the resulting model.WindowSize.
// An empty argument leaves the value unset.
func (s *WindowSize) Set(arg string) error {
	if arg == "" {
		return nil
	}
	sep := strings.IndexAny(arg, "x,")
	if sep < 0 {
		return fmt.Errorf("invalid window size %q: expected WxH", arg)
	}
	width, err := strconv.Atoi(strings.TrimSpace(arg[:sep]))
	if err != nil {
		return fmt.Errorf("invalid window size %q: %w", arg, err)
	}
	height, err := strconv.Atoi(strings.TrimSpace(arg[sep+1:]))
	if err != nil {
		return fmt.Errorf("invalid window size %q: %w", arg, err)
	}
	s.WindowSize = model.WindowSize{Width: width, Height: height}
	return nil
}

func (s *WindowSize) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

var _ flag.Getter = (*WindowSize)(nil)
