package compiler

import (
	"reflect"
	"testing"

	"github.com/0xddy/sing-box-acp-support/internal/topology"
)

func TestCompileRouteRuleSetUsesHTTPClientDetour(t *testing.T) {
	compiled, err := compileRouteRuleSet(topology.RouteRuleSet{
		Type:           "remote",
		Tag:            "remote-rules",
		Format:         "binary",
		URL:            "https://example.com/rules.srs",
		DownloadDetour: "direct",
	})
	if err != nil {
		t.Fatal(err)
	}

	wantHTTPClient := map[string]any{"detour": "direct"}
	if !reflect.DeepEqual(compiled["http_client"], wantHTTPClient) {
		t.Fatalf("http_client = %#v, want %#v", compiled["http_client"], wantHTTPClient)
	}
	if _, exists := compiled["download_detour"]; exists {
		t.Fatalf("deprecated download_detour should not be emitted: %#v", compiled)
	}
}
