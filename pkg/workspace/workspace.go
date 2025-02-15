package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"
)

type Workspace struct {
	Directory string
	Allocator *automation.Allocator

	stateCache *workspaceState
}

type workspaceState struct {
	Automations     []*automation.Automation
	didImplicitLoad bool
}

func (w *Workspace) Load(files ...string) error {
	err := w.load()
	if err != nil {
		return err
	}

	root := os.DirFS(w.Directory)
	p := config.NewParser(root)
	res := []*config.File{}

	for _, path := range files {
		file, diag := p.LoadConfigFile(path)
		if diag.HasErrors() {
			return diag
		}
		res = append(res, file)
	}
	w.state().appendFiles(res)
	return nil
}

func (w *Workspace) Automations() []*automation.Automation {
	err := w.load()
	if err != nil {
		logError(err)
	}
	return w.state().Automations
}

// Automation retrieves the automation by name
func (w *Workspace) Automation(name string) *automation.Automation {
	for _, auto := range w.Automations() {
		if auto.Name == name {
			return auto
		}
	}
	return nil
}

// Dir gets the workspace directory, normalized
func (w *Workspace) Dir() string {
	return w.actualDirectory()
}

// Autogun gets the directory where workspace metadata is stored
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

func (w *Workspace) load() error {
	return w.loadFilesFromWS()
}

func (w *Workspace) loadFilesFromWS() error {
	if w.state().didImplicitLoad {
		return nil
	}
	w.state().didImplicitLoad = true

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
			file, diag := p.LoadConfigFile(path)
			if diag.HasErrors() {
				return diag
			}
			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		return err
	}

	w.state().appendFiles(files)
	return nil
}

func (w *Workspace) state() *workspaceState {
	if w.stateCache == nil {
		w.stateCache = new(workspaceState)
	}
	return w.stateCache
}

// TODO This should not be API
func (w *Workspace) EnsureAllocator() *automation.Allocator {
	if w.Allocator == nil {
		w.Allocator = &automation.Allocator{}
	}
	return w.Allocator
}

func (s *workspaceState) appendFiles(files []*config.File) {
	for _, file := range files {
		for _, auto := range file.Automations {
			a, err := automation.Bind(auto)
			if err != nil {
				logError(err)
			}
			s.Automations = append(s.Automations, a)
		}
	}
}

func logError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
