package controller

import (
	"context"
	"testing"

	"github.com/kandev/kandev/internal/agent/registry"
	"github.com/kandev/kandev/internal/agent/settings/models"
	"github.com/kandev/kandev/internal/auth/authn"
	"github.com/kandev/kandev/internal/common/logger"
)

func TestAgentCatalogAuthority(t *testing.T) {
	st := newFakeStore()
	agent := &models.Agent{ID: "agent", Name: "catalog-agent"}
	st.agents[agent.ID], st.byName[agent.Name] = agent, agent
	st.profiles[agent.ID] = []*models.AgentProfile{{ID: "profile", AgentID: agent.ID, Name: "default", Enabled: true, CommandPrefix: "private-launcher"}}
	log, err := logger.NewLogger(logger.LoggingConfig{Level: "error", Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	ctrl := &Controller{repo: st, agentRegistry: registry.NewRegistry(log), logger: log}
	for _, tc := range []struct {
		name string
		id   authn.Identity
		full bool
	}{
		{"absent", authn.Identity{}, false},
		{"synthetic", authn.Identity{UserID: "local", Role: authn.RoleAdmin, Synthetic: true}, false},
		{"agent-owner", authn.Identity{UserID: "owner", Role: authn.RoleAdmin}, false},
		{"PAT", authn.Identity{UserID: "owner", Role: authn.RoleAdmin, TokenID: "pat"}, false},
		{"human", authn.Identity{UserID: "owner", Role: authn.RoleAdmin, SessionID: "session"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := authn.WithIdentity(context.Background(), tc.id)
			one, err := ctrl.GetAgent(ctx, agent.ID)
			if err != nil {
				t.Fatal(err)
			}
			list, err := ctrl.ListAgents(ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, got := range []string{one.Profiles[0].CommandPrefix, list.Agents[0].Profiles[0].CommandPrefix} {
				if (got == "private-launcher") != tc.full {
					t.Fatalf("wrong catalog authority: prefix=%q", got)
				}
			}
		})
	}
}
