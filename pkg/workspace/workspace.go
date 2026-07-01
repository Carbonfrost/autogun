// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/Carbonfrost/autogun/pkg/internal/contextkey"
	"github.com/Carbonfrost/autogun/pkg/model"
	cli "github.com/Carbonfrost/joe-cli"
	joeconfig "github.com/Carbonfrost/joe-cli/extensions/config"
)

type Workspace struct {
	cli.Action

	*joeconfig.Workspace

	Directory string
	Allocator *automation.Allocator

	model   *model.Model
	loadErr error
}

// New creates a new workspace
func New() *Workspace {
	ws := &Workspace{
		Workspace: joeconfig.NewWorkspace(),
	}
	ws.Action = cli.Pipeline(
		ws.Workspace.Action,
		ContextValue(ws),
	)
	return ws
}

// FromContext gets the Workspace from the context otherwise panics
func FromContext(ctx context.Context) *Workspace {
	return contextkey.Resolve(ctx, contextkey.Workspace).(*Workspace)
}

// ContextValue provides an action that sets the given value into the context.
// The only supported type is *Workspace.
func ContextValue(v *Workspace) cli.Action {
	return cli.WithContextValue(contextkey.Workspace, v)
}

func (w *Workspace) Pipeline() cli.Action {
	return w.Action
}

// Load scans the workspace for configuration files and builds a Model from
// the automations they declare.
func (w *Workspace) Load() (*model.Model, error) {
	files, err := w.loadFiles()
	if err != nil {
		return nil, err
	}

	m := model.New(files...)
	return m, nil
}

// Model obtains the model for the workspace. This method implicitly
// loads the workspace and panics when the workspace cannot be loaded.
// Investigate [Workspace.Load] to load while handling load errors
func (w *Workspace) Model() *model.Model {
	if w.model == nil && w.loadErr == nil {
		w.model, w.loadErr = w.Load()
	}
	if w.loadErr != nil {
		panic(w.loadErr)
	}
	return w.model
}

// Dir gets the workspace directory, normalized
func (w *Workspace) Dir() string {
	return w.actualDirectory()
}

// AutogunDir gets the directory where workspace metadata is stored
func (w *Workspace) AutogunDir() string {
	return filepath.Join(w.actualDirectory(), ".autogun")
}

func (w *Workspace) actualDirectory() string {
	if w.Directory == "" {
		res, _ := os.Getwd()
		return res
	}
	return w.Directory
}

func (w *Workspace) loadFiles() ([]*config.File, error) {
	root := os.DirFS(w.AutogunDir())
	p := config.NewParser(root)
	files := []*config.File{}
	err := fs.WalkDir(root, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".autog") ||
			strings.HasSuffix(path, ".autogun") ||
			strings.HasSuffix(path, ".hcl") ||
			strings.HasSuffix(path, ".autog.json") ||
			strings.HasSuffix(path, ".autogun.json") ||
			strings.HasSuffix(path, ".hcl.json") {
			file, diag := p.LoadFile(path)
			if diag.HasErrors() {
				return diag
			}
			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// TODO This should not be API
func (w *Workspace) EnsureAllocator() *automation.Allocator {
	if w.Allocator == nil {
		w.Allocator = &automation.Allocator{}
	}
	return w.Allocator
}

func logError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
