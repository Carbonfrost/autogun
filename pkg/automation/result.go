// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Result struct {
	Outputs     map[string]*json.RawMessage
	OutputFiles map[string]*[]byte
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
