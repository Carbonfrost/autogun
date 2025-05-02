// Copyright 2022 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"os"

	"github.com/Carbonfrost/autogun/pkg/cmd/autogun"
)

func main() {
	autogun.Run(os.Args)
}
