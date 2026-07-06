// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
)

type Protocol int

const (
	Chromedp Protocol = iota
)

func (e *Protocol) Value() automation.SupportedProtocol {
	switch *e {
	case Chromedp:
		return automation.ProtocolChromedp
	}
	return 0
}

func (*Protocol) Synopsis() string {
	return "{chromedp}"
}

func (e *Protocol) Set(arg string) error {
	switch strings.ToLower(arg) {
	case "chromedp":
		*e = Chromedp
	default:
		return fmt.Errorf("invalid value: %q", arg)
	}
	return nil
}

func (e Protocol) String() string {
	switch e {
	case Chromedp:
		return "chromedp"
	default:
	}
	return ""
}

var _ flag.Value = (*Protocol)(nil)
