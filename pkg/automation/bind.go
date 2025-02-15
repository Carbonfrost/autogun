package automation

import (
	"github.com/Carbonfrost/autogun/pkg/config"
)

// Bind generates the automation from configuration
func Bind(cfg *config.Automation) (*Automation, error) {
	return &Automation{
		Tasks: bindAutomation(cfg),
	}, nil
}
