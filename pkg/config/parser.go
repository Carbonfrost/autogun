// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package config

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Parser struct {
	fs fs.FS
	p  *hclparse.Parser
}

const (
	badIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."
)

var (
	identifierPattern = regexp.MustCompile("(?i)[_a-z][a-z0-9_-]*")
)

func NewParser(fs fs.FS) *Parser {
	if fs == nil {
		fs = os.DirFS(".")
	}

	return &Parser{
		fs: fs,
		p:  hclparse.NewParser(),
	}
}

func (p *Parser) LoadHCLFile(path string) (hcl.Body, hcl.Diagnostics) {
	src, err := fs.ReadFile(p.fs, path)

	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The file %q could not be read.", path),
			},
		}
	}

	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)

	switch {
	case strings.HasSuffix(path, ".json"):
		file, diags = p.p.ParseJSON(src, path)
	default:
		file, diags = p.p.ParseHCL(src, path)
	}

	if file == nil || file.Body == nil {
		return hcl.EmptyBody(), diags
	}

	return file.Body, diags
}

func (p *Parser) LoadConfigFile(path string) (*File, hcl.Diagnostics) {
	body, diags := p.LoadHCLFile(path)
	if body == nil {
		return nil, diags
	}

	return decodeFile(path, body)
}

func tryLabel(b *hcl.Block, n int) string {
	if n < len(b.Labels) {
		return b.Labels[n]
	}
	return ""
}

func tryLabelRange(b *hcl.Block, n int) (res hcl.Range) {
	if n < len(b.LabelRanges) {
		res = b.LabelRanges[n]
	}
	return
}

func validIdentifier(name string) bool {
	return len(name) == 0 || identifierPattern.MatchString(name)
}
