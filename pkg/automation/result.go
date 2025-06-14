// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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

func Execute(ctx context.Context, allocator *Allocator, a *Automation) (*Result, error) {
	res := newResult()
	ctx, cancel := allocator.newContext(
		withAutomationResult(ctx, res),
	)
	defer cancel()

	var emulate Task = TaskFunc(nil)
	if dev, ok := allocator.resolveDevice(); ok {
		fmt.Fprintf(os.Stderr, "Emulating device %s (%s)\n", dev.Device().Name, allocator.DeviceID)
		emulate = chromedp.Emulate(dev)
	}

	return res, chromedp.Run(ctx, emulate, a)
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
