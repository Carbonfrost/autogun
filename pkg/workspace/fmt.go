// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Carbonfrost/autogun/pkg/config/format"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
)

type FormatParams struct {
	Files *cli.FileSet
	Write bool
}

// Fmt returns an action which formats source. The parameters
// specify which files to format. When the parameters are not specified,
// the parameters are gathered from the action's context. Also when
// the parameters are not specified, the action contains an initializer
// that contains args and flags meant to provide useful default configuration.
func Fmt(paramsopt ...FormatParams) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "fmt",
			HelpText: "Apply formatting to files in the workspace",
		},
		bind.Call(formatSpec, useFormatParams()),
	)
}

func useFormatParams() bind.ActionBinder[FormatParams] {
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
			cli.AddFlags([]*cli.Flag{
				{
					Name:     "write",
					Aliases:  []string{"w"},
					Value:    new(bool),
					HelpText: "update source files in-place instead of writing to stdout",
				},
			}...),
		),
		bind.Func[FormatParams](func(c *cli.Context) (FormatParams, error) {
			return FormatParams{
				Files: c.FileSet("files"),
				Write: c.Bool("write"),
			}, nil
		}),
	)
}

func formatSpec(f FormatParams) error {
	overwrite := f.Write

	files, err := enumerateFiles(f.Files)
	if err != nil {
		return err
	}

	return processFiles(overwrite, files)
}

func processFiles(overwrite bool, files []string) error {
	if len(files) == 0 {
		if overwrite {
			return errors.New("error: cannot use -w without source filenames")
		}

		_, err := processFile("<stdin>", os.Stdin, false)
		return err
	}

	var anyChanged bool
	for _, path := range files {
		switch dir, err := os.Stat(path); {
		case err != nil:
			return err
		case dir.IsDir():
			return fmt.Errorf("can't format directory %s", path)
		default:
			changed, err := processFile(path, nil, overwrite)
			if changed {
				fmt.Fprintf(os.Stderr, "%s\n", path)
				anyChanged = true
			}
			if err != nil {
				return err
			}
		}
	}

	if anyChanged {
		return cli.Exit(1)
	}

	return nil
}

func processFile(fn string, in *os.File, overwrite bool) (changed bool, err error) {
	if in == nil {
		in, err = os.Open(fn)
		if err != nil {
			return false, fmt.Errorf("failed to open %s: %s", fn, err)
		}
	}

	inSrc, err := io.ReadAll(in)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %s", fn, err)
	}

	outSrc := format.Source(inSrc)
	changed = !bytes.Equal(inSrc, outSrc)

	if overwrite {
		if changed {
			return true, os.WriteFile(fn, outSrc, 0644)
		}

		return false, nil
	}

	_, err = os.Stdout.Write(outSrc)
	return false, err
}
