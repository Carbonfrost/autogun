// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Source applies source formatting to the given file
func Source(inSrc []byte) []byte {
	return hclwrite.Format(inSrc)
}
