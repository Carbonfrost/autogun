// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"golang.org/x/term"
)

type CheckParams struct {
	Files *cli.FileSet
}

func Check() cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "check",
			HelpText: "Parse files to look for syntax errors",
		},
		bind.Call(checkSpec, useCheckParams()),
	)
}

func useCheckParams() bind.ActionBinder[CheckParams] {
	return bind.NewActionBinder(
		cli.Pipeline(
			cli.AddArgs([]*cli.Arg{
				{
					Name:  "files",
					Value: new(cli.FileSet),
					NArg:  cli.TakeUntilNextFlag,
					Uses:  cli.Accessory("recursive", (*cli.FileSet).RecursiveFlag, cli.HelpText("format directories recursively")),
				},
			}...),
		),
		bind.Func[CheckParams](func(c *cli.Context) (CheckParams, error) {
			return CheckParams{
				Files: c.FileSet("files"),
			}, nil
		}),
	)
}

func checkSpec(c CheckParams) error {
	var anyErrors bool
	parser := hclparse.NewParser()
	color := term.IsTerminal(int(os.Stderr.Fd()))
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}
	diagWriter := hcl.NewDiagnosticTextWriter(os.Stderr, parser.Files(), uint(w), color)
	files, err := enumerateFiles(c.Files)
	if err != nil {
		return err
	}

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return unableToCheck(path, err)
		}

		_, diags := parser.ParseHCL(data, path)
		diagWriter.WriteDiagnostics(diags)

		if diags.HasErrors() {
			anyErrors = true
		}
	}
	if anyErrors {
		return cli.Exit("one or more files contained errors", 1)
	}
	return nil
}

func unableToCheck(path string, err error) error {
	return fmt.Errorf("unable to check %s: %w", path, err)
}

func enumerateFiles(fs *cli.FileSet) ([]string, error) {
	var files []string
	for input, err := range fs.All() {
		if err != nil {
			return nil, err
		}

		s, _ := input.Stat()
		if s.IsDir() {
			if fs.Recursive {
				continue
			}
			return nil, fmt.Errorf("can't process directory: %s", input.Name)
		}

		if !detectFile(input.Name) {
			if fs.Recursive {
				continue
			}
			return nil, fmt.Errorf("unsupported file for process: %s", input.Name)
		}
		files = append(files, input.Name)
	}
	return files, nil
}

func detectFile(path string) bool {
	return strings.HasSuffix(path, ".autog") ||
		strings.HasSuffix(path, ".autogun") ||
		strings.HasSuffix(path, ".hcl") ||
		strings.HasSuffix(path, ".autog.json") ||
		strings.HasSuffix(path, ".autogun.json") ||
		strings.HasSuffix(path, ".hcl.json")
}
