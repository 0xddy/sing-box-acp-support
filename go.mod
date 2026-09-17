module github.com/0xddy/sing-box-acp-support

go 1.25.5

require (
	github.com/acp/node-agent v0.0.0
	github.com/miekg/dns v1.1.72
	github.com/sagernet/sing v0.9.4
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba
	golang.org/x/mod v0.38.0
	golang.org/x/text v0.40.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260715232425-e75dac1f907d // indirect
	google.golang.org/grpc v1.82.1 // indirect
)

replace github.com/acp/node-agent => ../node-agent/src
