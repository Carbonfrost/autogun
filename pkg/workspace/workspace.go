package workspace

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/chromedp/chromedp"
)

const (
	workspaceKey contextKey = "workspaceKey"
)

type Workspace struct {
	Directory string
	Allocator *Allocator

	stateCache *workspaceState
}

type workspaceState struct {
	Automations     []*config.Automation
	didImplicitLoad bool
}

// Workspace gets the Workspace from the context
func FromContext(ctx context.Context) *Workspace {
	return ctx.Value(workspaceKey).(*Workspace)
}

// SetContextWorkspace will set the Workspace
func SetContextWorkspace(ctx context.Context, ws *Workspace) context.Context {
	return context.WithValue(ctx, workspaceKey, ws)
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

func (w *Workspace) Automations() []*config.Automation {
	return w.state().Automations
}

func (w *Workspace) automation(name string) *config.Automation {
	for _, auto := range w.state().Automations {
		if auto.Name == name {
			return auto
		}
	}
	return nil
}

func (w *Workspace) Execute(automation string) (*AutomationResult, error) {
	err := w.load()
	if err != nil {
		return nil, err
	}

	auto := w.automation(automation)
	if auto == nil {
		return nil, fmt.Errorf("automation not found %q", automation)
	}
	return w.executeCore(auto)
}

func (w *Workspace) executeCore(auto *config.Automation) (*AutomationResult, error) {
	res := NewAutomationResult()
	tasks := res.bindAutomation(auto)
	ctx, cancel := w.ensureAllocator().newContext(context.Background())
	defer cancel()

	return res, chromedp.Run(ctx, tasks...)
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
			strings.HasSuffix(path, ".hcl") ||
			strings.HasSuffix(path, ".autog.json") ||
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

func (w *Workspace) ensureAllocator() *Allocator {
	if w.Allocator == nil {
		w.Allocator = &Allocator{}
	}
	return w.Allocator
}

func (s *workspaceState) appendFiles(files []*config.File) {
	for _, file := range files {
		s.Automations = append(s.Automations, file.Automations...)
	}
}
