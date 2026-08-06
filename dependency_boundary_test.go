package support

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestProductionDependencyBoundary(t *testing.T) {
	moduleFile, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if bytes.Contains(moduleFile, []byte("github.com/sagernet/sing-box")) {
		t.Fatal("root go.mod must not require or replace github.com/sagernet/sing-box")
	}

	command := exec.Command("go", "list", "-mod=readonly", "-deps", "./...")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("list production dependencies: %v\n%s", err, output)
	}

	forbiddenPrefixes := []string{
		"github.com/sagernet/sing-box",
		"github.com/sagernet/sing-quic",
		"github.com/sagernet/quic-go",
		"github.com/sagernet/sing-tun",
		"github.com/sagernet/gvisor",
		"github.com/quic-go/",
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		dependency := strings.TrimSpace(scanner.Text())
		for _, prefix := range forbiddenPrefixes {
			if dependency == prefix || strings.HasPrefix(dependency, prefix+"/") {
				t.Errorf("forbidden production dependency: %s", dependency)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan dependencies: %v", err)
	}
}
