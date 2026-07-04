// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"runtime/debug"
	"strings"
)

var Version VersionInfo

type VersionInfo struct {
	Version string
	Modules map[string]string
}

var keyModules = map[string]bool{
	"cdproto":  true,
	"chromedp": true,
}

func init() {
	info, _ := debug.ReadBuildInfo()

	res := map[string]string{}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, d := range info.Deps {
			index := strings.LastIndex(d.Path, "/")
			if index > 0 {
				name := d.Path[index+1:]
				if keyModules[name] {
					res[name] = d.Version
				}
			}
		}
	}
	Version = VersionInfo{
		Modules: res,
		Version: info.Main.Version,
	}

}
