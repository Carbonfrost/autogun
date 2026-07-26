// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Network holds network filtering configuration for an automation.
type Network struct {
	DeclRange hcl.Range
	Filters   []*NetworkFilter
}

// NetworkFilter specifies criteria for matching network responses and what to do
// when a match is found.
type NetworkFilter struct {
	DeclRange hcl.Range
	Type      string                 // optional ResourceType name
	Response  *NetworkFilterResponse // optional response URL matching
	Download  *NetworkDownload       // optional; when present, save body to OutputFiles
	Log       *NetworkLog            // optional; when present, log the match
}

// NetworkFilterResponse holds URL glob patterns for matching response URLs.
type NetworkFilterResponse struct {
	DeclRange hcl.Range
	URLs      []string
}

// NetworkDownload signals that a matched response body should be saved.
type NetworkDownload struct {
	DeclRange hcl.Range
}

// NetworkLog signals that a matched response should be logged.
type NetworkLog struct {
	DeclRange hcl.Range
}

var (
	networkBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "filter"},
		},
	}

	networkFilterBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "response"},
			{Type: "download"},
			{Type: "log"},
		},
	}

	networkFilterResponseBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "url"},
		},
	}

	networkDownloadBlockSchema = emptyBlockSchema
	networkLogBlockSchema      = emptyBlockSchema

	mappingNetworkFilterBlocks = blockMapping[*NetworkFilter]{
		"filter": decodeNetworkFilterBlock,
	}
)

// automationNetworkSchema is used in a separate PartialContent pass so that
// the network block is not mixed into the task-block mapping.
var automationNetworkSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "network"},
	},
}

func decodeNetworkBlock(block *hcl.Block) (*Network, hcl.Diagnostics) {
	n := new(Network)
	return reduce(
		n,
		block,
		supportsDeclRange(&n.DeclRange),
		supportsPartialContentSchema(
			networkBlockSchema,
			appendsTo(&n.Filters, mappingNetworkFilterBlocks),
		),
	)
}

func decodeNetworkFilterBlock(block *hcl.Block) (*NetworkFilter, hcl.Diagnostics) {
	f := new(NetworkFilter)
	return reduce(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			networkFilterBlockSchema,
			withAttribute("type", &f.Type),
			decodeNetworkFilterSubBlocks(f),
		),
	)
}

func decodeNetworkFilterSubBlocks(f *NetworkFilter) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		var diags hcl.Diagnostics
		for _, block := range content.Blocks {
			switch block.Type {
			case "response":
				r, rDiags := decodeNetworkFilterResponseBlock(block)
				diags = append(diags, rDiags...)
				if r != nil {
					f.Response = r
				}
			case "download":
				d, dDiags := decodeNetworkDownloadBlock(block)
				diags = append(diags, dDiags...)
				if d != nil {
					f.Download = d
				}
			case "log":
				l, lDiags := decodeNetworkLogBlock(block)
				diags = append(diags, lDiags...)
				if l != nil {
					f.Log = l
				}
			}
		}
		return diags
	}
}

func decodeNetworkFilterResponseBlock(block *hcl.Block) (*NetworkFilterResponse, hcl.Diagnostics) {
	r := new(NetworkFilterResponse)
	return reduce(
		r,
		block,
		supportsDeclRange(&r.DeclRange),
		supportsPartialContentSchema(
			networkFilterResponseBlockSchema,
			withAttribute("url", &r.URLs),
		),
	)
}

func decodeNetworkDownloadBlock(block *hcl.Block) (*NetworkDownload, hcl.Diagnostics) {
	d := new(NetworkDownload)
	return reduce(
		d,
		block,
		supportsDeclRange(&d.DeclRange),
		supportsPartialContentSchema(networkDownloadBlockSchema),
	)
}

func decodeNetworkLogBlock(block *hcl.Block) (*NetworkLog, hcl.Diagnostics) {
	l := new(NetworkLog)
	return reduce(
		l,
		block,
		supportsDeclRange(&l.DeclRange),
		supportsPartialContentSchema(networkLogBlockSchema),
	)
}

func withNetworkBlock(target **Network) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		var diags hcl.Diagnostics
		for _, block := range content.Blocks {
			if block.Type != "network" {
				continue
			}
			n, nDiags := decodeNetworkBlock(block)
			diags = append(diags, nDiags...)
			if n != nil {
				*target = n
			}
		}
		return diags
	}
}
