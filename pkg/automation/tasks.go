// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"os"
)

// Print returns a Task that writes its operands to stderr, formatting them in
// the manner of fmt.Print and appending a newline.
func Print(args ...any) Task {
	return TaskFunc(func(context.Context) error {
		fmt.Fprintln(os.Stderr, args...)
		return nil
	})
}

// Printf returns a Task that writes a formatted message to stderr, formatting
// it in the manner of fmt.Printf and appending a newline.
func Printf(format string, args ...any) Task {
	return TaskFunc(func(context.Context) error {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintln(os.Stderr)
		return nil
	})
}

func taskThunk(factory func(context.Context) (Task, error)) Task {
	return TaskFunc(func(c context.Context) error {
		auto, err := factory(c)
		if err != nil {
			return err
		}
		return auto.Do(c)
	})
}
