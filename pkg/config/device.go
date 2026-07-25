// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Device describes an inline device definition used to initialize all
// attributes of a device to emulate.
type Device struct {
	DeclRange hcl.Range
	Name      string
	UserAgent string
	Width     int64
	Height    int64
	Scale     float64
	Landscape bool
	Mobile    bool
	Touch     bool
}

var (
	deviceBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "name"},
			{Name: "user_agent"},
			{Name: "width"},
			{Name: "height"},
			{Name: "scale"},
			{Name: "landscape"},
			{Name: "mobile"},
			{Name: "touch"},
		},
	}
)

func decodeDeviceBlock(block *hcl.Block) (*Device, hcl.Diagnostics) {
	f := new(Device)
	return reduce(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			deviceBlockSchema,
			withAttribute("name", &f.Name),
			withAttribute("user_agent", &f.UserAgent),
			withAttribute("width", &f.Width),
			withAttribute("height", &f.Height),
			withAttribute("scale", &f.Scale),
			withAttribute("landscape", &f.Landscape),
			withAttribute("mobile", &f.Mobile),
			withAttribute("touch", &f.Touch),
		),
	)
}

func supportsDeviceBlock(dev **Device) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		var diags hcl.Diagnostics

		for _, block := range content.Blocks {
			if block.Type != "device" {
				continue
			}
			cfg, cfgDiags := decodeDeviceBlock(block)
			if cfg != nil {
				*dev = cfg
			}
			diags = append(diags, cfgDiags...)
		}
		return diags
	}
}
