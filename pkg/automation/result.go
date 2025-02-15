package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chromedp/chromedp"
)

type Result struct {
	Outputs     map[string]*json.RawMessage
	OutputFiles map[string]*[]byte
}

func Execute(allocator *Allocator, a *Automation) (*Result, error) {
	res := newResult()
	ctx, cancel := allocator.newContext(
		withAutomationResult(context.Background(), res),
	)
	defer cancel()

	return res, chromedp.Run(ctx, a)
}

func newResult() *Result {
	return &Result{
		Outputs:     map[string]*json.RawMessage{},
		OutputFiles: map[string]*[]byte{},
	}
}

// TODO This should not necessarily be API
func (r *Result) PersistOutputFiles() {
	for name, f := range r.OutputFiles {
		file, err := os.Create(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error persisting output files: %v\n", err)
			continue
		}

		_, err = io.Copy(file, bytes.NewReader(*f))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error persisting output files: %v\n", err)
			continue
		}
	}
}

func mustAutomationResult(c context.Context) *Result {
	return c.Value(automationResultKey).(*Result)
}

func withAutomationResult(c context.Context, ar *Result) context.Context {
	return context.WithValue(c, automationResultKey, ar)
}
