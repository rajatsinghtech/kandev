package handlers

import (
	"context"
	"testing"

	ws "github.com/kandev/kandev/pkg/websocket"
)

func TestMCPDoesNotGrantLaunchConfigurationAuthority(t *testing.T) {
	for _, action := range []string{ws.ActionMCPUpdateAgent, ws.ActionMCPCreateAgentProfile,
		ws.ActionMCPUpdateAgentProfile, ws.ActionMCPDeleteAgentProfile, ws.ActionMCPGetSettings,
		ws.ActionMCPUpdateSettings, ws.ActionMCPGetMcpConfig, ws.ActionMCPUpdateMcpConfig,
		ws.ActionMCPCreateExecutorProfile, ws.ActionMCPUpdateSharedPrompt} {
		t.Run(action, func(t *testing.T) {
			d := &guardedMCPDispatcher{Dispatcher: ws.NewDispatcher(), handlers: &Handlers{}}
			called := false
			d.RegisterFunc(action, func(ctx context.Context, msg *ws.Message) (*ws.Message, error) {
				called = true
				return ws.NewResponse(msg.ID, msg.Action, "sensitive-config")
			})
			response, err := d.Dispatch(context.Background(), &ws.Message{ID: "request", Action: action})
			if err != nil || response == nil || called {
				t.Fatalf("MCP configuration authority granted: called=%v err=%v", called, err)
			}
		})
	}
}
