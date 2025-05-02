// Copyright 2022 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package autogun_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAutogun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autogun Suite")
}
