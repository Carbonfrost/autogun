package autogun

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/contextual"
	cli "github.com/Carbonfrost/joe-cli"
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

		devices[id] = i.Device()
	}
}

func FlagsAndArgs() cli.Action {
	return cli.Pipeline(
		cli.AddFlags([]*cli.Flag{
			{Uses: SetBrowserURL()},
		}...),
	)
}

func RunAutomation() cli.Action {
	return cli.Setup{
		Action: func(c *cli.Context) error {
			ws := contextual.Workspace(c)
			err := ws.Load(c.FileSet("files").Files...)
			if err != nil {
				return err
			}

			for _, auto := range ws.Automations() {
				res, err := ws.ExecuteCore(auto)
				if err != nil {
					return err
				}

				data, _ := json.MarshalIndent(res.Outputs, "", "    ")
				os.Stdout.Write(data)

				res.PersistOutputFiles()
			}
			return nil
		},
	}
}

func ListDevices() cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "list-devices",
			HelpText: "list available devices to emulate",
			Options:  cli.Exits | cli.NonPersistent,
			Value:    new(bool),
		},
		cli.At(cli.ActionTiming, cli.ActionOf(func() {
			for _, s := range slices.Sorted(maps.Keys(devices)) {
				fmt.Println(s)
			}
		})))
}

func SetBrowserURL(v ...string) cli.Action {
	return cli.Pipeline(
		&cli.Prototype{
			Name:     "browser",
			Aliases:  []string{"b"},
			HelpText: "Connect to the running browser instance by {URL}",
		},
		withBinding((*automation.Allocator).SetBrowserURL, v...),
	)
}

func withBinding[V any](binder func(*automation.Allocator, V) error, args ...V) cli.Action {
	return cli.BindContext(allocatorFromContext, binder, args...)
}

func allocatorFromContext(c context.Context) *automation.Allocator {
	return contextual.Workspace(c).EnsureAllocator()
}
