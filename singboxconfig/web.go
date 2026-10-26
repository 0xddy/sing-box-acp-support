package singboxconfig

import (
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	"github.com/acp/node-agent/api/provider"
)

// ValidateWeb validates the shared website contract using the mirrored core options.
func ValidateWeb(config provider.WebConfig) error { return (option.WebOptions(config)).Validate() }
