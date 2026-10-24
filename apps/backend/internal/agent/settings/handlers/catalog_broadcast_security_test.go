package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kandev/kandev/internal/agent/settings/dto"
	ws "github.com/kandev/kandev/pkg/websocket"
)

func TestAgentCatalogBroadcastOmitsLaunchAuthority(t *testing.T) {
	for _, workspace := range []string{"", "workspace"} {
		t.Run("workspace="+workspace, func(t *testing.T) {
			hub := newWorkspaceDuplicateHub()
			h := &Handlers{hub: hub}
			profile := &dto.AgentProfileDTO{ID: "profile", AgentID: "agent", WorkspaceID: workspace,
				Name: "selection", CommandPrefix: "private-launcher", ProviderBaseURL: "https://private-provider"}
			h.broadcastProfileEvent(context.Background(), ws.ActionAgentProfileDeleted, profile)
			messages := hub.global
			if workspace != "" {
				messages = hub.workspaceMsgs[workspace]
			}
			if len(messages) != 1 {
				t.Fatalf("selection notification missing: %d", len(messages))
			}
			payload := string(messages[0].Payload)
			if strings.Contains(payload, "private-") {
				t.Fatal("broadcast disclosed launcher configuration")
			}
			var result struct {
				Profile dto.AgentProfileDTO `json:"profile"`
			}
			if err := json.Unmarshal(messages[0].Payload, &result); err != nil || result.Profile.ID != profile.ID {
				t.Fatalf("redacted selection lost identity: %v", err)
			}
			if profile.CommandPrefix != "private-launcher" {
				t.Fatal("notification redaction mutated stored configuration")
			}
		})
	}
}
