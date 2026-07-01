// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build json_marshal

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Carbonfrost/autogun/pkg/cmd/autogun"
	"github.com/Carbonfrost/joe-cli/extensions/marshal"
)

// This prints out marshal information about the app
// go run -tags json_marshal ./cmd/autogun

func main() {
	app := autogun.NewApp()
	_, err := app.Initialize(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	m := marshal.From(app)
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "    ")
	e.Encode(m)
}
