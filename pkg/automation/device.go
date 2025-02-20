package automation

import (
	"maps"
	"slices"
	"strings"

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

func DeviceIDs() []string {
	return slices.Sorted(maps.Keys(devices))
}
