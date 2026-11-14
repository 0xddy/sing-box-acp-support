package configtarget

import (
	C "github.com/0xddy/sing-box-acp-support/internal/singbox/constant"
	"github.com/0xddy/sing-box-acp-support/internal/singbox/option"
	"github.com/acp/node-agent/pkg/configcompile/provider"
)

type Schema struct{}

func (Schema) ValidateWeb(web provider.WebConfig) error { return (option.WebOptions(web)).Validate() }
func (Schema) ValidNetworkStrategy(value string) bool {
	_, ok := C.StringToNetworkStrategy[value]
	return ok
}
func (Schema) ServerIsDomain(server string) bool {
	return (option.DNSServerAddressOptions{Server: server}).ServerIsDomain()
}
