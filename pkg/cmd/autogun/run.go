package autogun

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/Carbonfrost/autogun/pkg/contextual"
	"github.com/Carbonfrost/autogun/pkg/workspace"
	cli "github.com/Carbonfrost/joe-cli"
	"github.com/chromedp/chromedp"
)

func RunAutomation(c *cli.Context) error {
	// Build a single automation given the source identified in the context,
	// which are loaded as files (or navigation URLs). The expression evaluation
	// functions append additional tasks.
	auto, err := convertSources(c)
	if err != nil {
		return err
	}

	exp := ensurePrinter(c.Expression("expression"))
	err = exp.Evaluate(c, auto)
	if err != nil {
		return err
	}

	ws := contextual.Workspace(c.Context())
	results, err := automation.Execute(c.Context(), ws.EnsureAllocator(), auto)
	if err != nil {
		return err
	}

	// Generate output from the output variables and persistent files
	data, _ := json.MarshalIndent(results.Outputs, "", "    ")
	os.Stdout.Write(data)
	results.PersistOutputFiles()

	return nil
}

func convertSources(c *cli.Context) (*automation.Automation, error) {
	var result []automation.Task
	for _, source := range c.List("sources") {
		if source == "." {
			continue

		} else if source == "-" {
			// TODO Support reading the automation from an input file
			return nil, fmt.Errorf("not yet implemented: read automation from stdin")

		} else if looksLikeURL(source) {
			result = append(result, chromedp.Navigate(source))

		} else {
			result = append(result, runSource(source))
		}
	}
	return &automation.Automation{
		Tasks: result,
	}, nil
}

func runSource(source string) automation.Task {
	return automation.TaskFunc(func(c context.Context) error {
		// TODO Dubious to parse the file like this without the
		// workspace loading first. Also loadOne assumes one automation
		// per file which may not be true

		// TODO Also dubious to perform this within the task (rather
		// than fail earlier), but it is convenient for implementing the
		// expression
		ws := contextual.Workspace(c)
		auto, err := loadOne(ws, source)

		if err != nil {
			return err
		}
		return auto.Do(c)
	})
}

func flow(name string) automation.Task {
	return automation.TaskFunc(func(c context.Context) error {
		a := contextual.Workspace(c).Automation(name)
		if a == nil {
			return fmt.Errorf("automation not found %q", name)
		}
		return a.Do(c)
	})
}

func looksLikeURL(addr string) bool {
	return strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")
}

func loadOne(w *workspace.Workspace, path string) (*automation.Automation, error) {
	root := os.DirFS(w.Directory)
	p := config.NewParser(root)
	file, diag := p.LoadConfigFile(path)
	if diag.HasErrors() {
		return nil, diag
	}

	return automation.Bind(file.Automations[0])
}
