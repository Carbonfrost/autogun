package autogun

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/Carbonfrost/joe-cli"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"golang.org/x/term"
)

func Check() *cli.Command {
	return &cli.Command{
		Name:     "check",
		HelpText: "Parse files to look for syntax errors",
		Args: []*cli.Arg{
			{
				Name:    "files",
				Value:   new(cli.FileSet),
				Options: cli.Merge,
				NArg:    cli.TakeUntilNextFlag,
			},
		},
		Flags: []*cli.Flag{
			{
				Name:     "recursive",
				Aliases:  []string{"r"},
				HelpText: "format directories recursively",
				Uses:     cli.BindIndirect("files", (*cli.FileSet).SetRecursive, true),
			},
		},
		Action: checkCommand,
	}
}

func checkCommand(c *cli.Context) error {
	var anyErrors bool
	parser := hclparse.NewParser()
	color := term.IsTerminal(int(os.Stderr.Fd()))
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}
	diagWriter := hcl.NewDiagnosticTextWriter(os.Stderr, parser.Files(), uint(w), color)
	files, err := enumerateFiles(c.FileSet("files"))
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
