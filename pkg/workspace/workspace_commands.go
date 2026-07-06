package workspace

import (
	"fmt"

	"github.com/Carbonfrost/autogun/pkg/model"
	"github.com/Carbonfrost/joe-cli"
	"github.com/Carbonfrost/joe-cli/extensions/bind"
)

// ListDevices provides an action which prints out the information about devices.
// If present, detailedopt controls the verbosity. When not present, a flag
// --verbose is consulted to control it
func ListDevices(detailedopt ...bool) cli.Action {
	var binder bind.Binder[bool]
	switch len(detailedopt) {
	case 1:
		binder = bind.Exact(detailedopt[0])
	case 0:
		binder = bind.Seen("verbose")
	default:
		panic("detailedopt must specify only one argument")
	}

	return cli.Pipeline(
		&cli.Prototype{
			Name:     "list-devices",
			HelpText: "list available devices to emulate",
			Options:  cli.Exits | cli.NonPersistent,
			Value:    new(bool),
		},
		bind.Call(printDevicesHelper, binder),
	)
}

func printDevicesHelper(detailed bool) error {
	var printer = func(s model.Device) {
		fmt.Print(s.ID, "\t", s.Name, "\n")
	}

	if detailed {
		fmt.Print("id", "\t", "name", "\t", "user agent", "\t", "width", "\t", "height", "\t",
			"scale", "\t", "landscape", "\t", "mobile", "\t", "touch", "\n")
		printer = func(d model.Device) {
			fmt.Print(d.ID, "\t", d.Name, "\t", d.UserAgent, "\t", d.Width, "\t", d.Height, "\t",
				d.Scale, "\t", d.Landscape, "\t", d.Mobile, "\t", d.Touch, "\n")
		}
	}

	for _, d := range model.Devices() {
		printer(d)
	}
	return nil
}
