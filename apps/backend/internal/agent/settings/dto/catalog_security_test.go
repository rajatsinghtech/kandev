package dto

import (
	"encoding/json"
	"strings"
	"testing"
)

// @covers AC-AUTH-OPERATOR-CONTROL-002.3
func TestAgentCatalogOmitsLaunchAuthority(t *testing.T) {
	source := AgentDTO{ID: "agent", Name: "Claude", MCPConfigPath: "private-mcp-path",
		TUIConfig: &TUIConfigDTO{Command: "private-command"},
		Profiles: []AgentProfileDTO{{ID: "profile", Name: "Default", Model: "model", Enabled: true,
			CommandPrefix: "private-prefix", ProviderBaseURL: "https://private-endpoint", ConfigOptions: map[string]string{"token": "private-option"}}}}
	catalog := RedactedAgentCatalog(source)
	data, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "private-") {
		t.Fatalf("catalog leaked launch authority: %s", data)
	}
	if catalog.ID != source.ID || len(catalog.Profiles) != 1 || catalog.Profiles[0].Model != "model" {
		t.Fatal("catalog lost selection metadata")
	}
	if source.Profiles[0].CommandPrefix != "private-prefix" {
		t.Fatal("redaction mutated internal configuration")
	}
}
