// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	cli "github.com/Carbonfrost/joe-cli"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func Fmt() *cli.Command {
	return &cli.Command{
		Name:     "fmt",
		HelpText: "Apply formatting to files in the workspace",
		Args: []*cli.Arg{
			{
				Name:  "files",
				Value: new(cli.FileSet),
				NArg:  cli.TakeUntilNextFlag,
			},
		},
		Flags: []*cli.Flag{
			{
				Name:     "recursive",
				Aliases:  []string{"r"},
				HelpText: "format directories recursively",
				Uses:     cli.BindIndirect("files", (*cli.FileSet).SetRecursive, true),
			},
			{
				Name:     "overwrite",
				Aliases:  []string{"w"},
				Value:    new(bool),
				HelpText: "update source files in-place instead of writing to stdout",
			},
		},
		Action: fmtCommand,
	}
}

func fmtCommand(c *cli.Context) error {
	overwrite := c.Bool("overwrite")

	files, err := enumerateFiles(c.FileSet("files"))
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

	outSrc := hclwrite.Format(inSrc)
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
