package automation

import (
	"github.com/Carbonfrost/autogun/pkg/config"
)

// Option provides an option to the automation builder
type Option func(*Automation)

// Binder provides the logic
type Binder func(*config.Automation) (*Automation, error)

// Bind converts the configuration into an automation
func Bind(cfg *config.Automation, using Binder, opts ...Option) (*Automation, error) {
	if using == nil {
		using = UsingChromedp
	}
	auto, err := using(cfg)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o(auto)
	}
	return auto, nil
}

// UsingChromedp is a Binder for using Chrome DevTools Protocol and headless
// Chrome/Chromium to run the automation
func UsingChromedp(cfg *config.Automation) (*Automation, error) {
	return &Automation{
		Name:  cfg.Name,
		Tasks: bindAutomation(cfg),
	}, nil
}

var (
	_ Binder = UsingChromedp
)
