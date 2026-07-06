// Copyright 2025, 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"cmp"
	"slices"
	"strings"
	"sync"

	"github.com/chromedp/chromedp/device"
)

var devices map[string]device.Info

func init() {
	devices = map[string]device.Info{}
	for i := device.Reset; i <= device.MotoG4landscape; i++ {
		id := strings.ReplaceAll(i.Device().Name, " ", "")
		id = strings.ReplaceAll(id, ")", "")
		id = strings.ReplaceAll(id, "(", "_")
		id = strings.ReplaceAll(id, "+", "plus")
		if id == "" {
			continue
		}

		devices[id] = i.Device()
	}
}

// Devices obtains all the devices
func Devices() []Device {
	return deviceList()
}

// LookupDevice looks up a device by ID
func LookupDevice(id string) (Device, bool) {
	d, ok := devices[id]
	return convertDevice(id, d), ok
}

// Device describes a device for emulation
type Device struct {
	// ID provides the ID of the device
	ID string

	// Name is the device name.
	Name string

	// UserAgent is the device user agent string.
	UserAgent string

	// Width is the viewport width.
	Width int64

	// Height is the viewport height.
	Height int64

	// Scale is the device viewport scale factor.
	Scale float64

	// Landscape indicates whether or not the device is in landscape mode or
	// not.
	Landscape bool

	// Mobile indicates whether it is a mobile device or not.
	Mobile bool

	// Touch indicates whether the device has touch enabled.
	Touch bool
}

var deviceList = sync.OnceValue(func() []Device {
	res := make([]Device, 0, len(devices))
	for id, d := range devices {
		res = append(res, convertDevice(id, d))
	}
	slices.SortFunc(res, deviceByID)
	return res
})

func deviceByID(x, y Device) int {
	return cmp.Compare(x.ID, y.ID)
}

func convertDevice(id string, d device.Info) Device {
	return Device{
		ID:        id,
		Name:      d.Name,
		UserAgent: d.UserAgent,
		Width:     d.Width,
		Height:    d.Height,
		Scale:     d.Scale,
		Landscape: d.Landscape,
		Mobile:    d.Mobile,
		Touch:     d.Touch,
	}
}
