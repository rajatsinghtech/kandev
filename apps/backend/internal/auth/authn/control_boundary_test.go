package authn

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// @covers AC-AUTH-ADMIN-CONTAINMENT-002.1
// @covers AC-AUTH-OPERATOR-CONTROL-002.2
func TestControlPlaneBoundary(t *testing.T) {
	admin := Identity{UserID: "a", OrgID: "org", Role: RoleAdmin, SessionID: "human"}
	operator := admin
	operator.Instance = true
	member := admin
	member.Role = RoleMember
	synthetic := Identity{UserID: "local", Role: RoleAdmin, Synthetic: true}
	pat := admin
	pat.SessionID, pat.TokenID = "", "automation"
	for _, tc := range []struct {
		name, method, path string
		identity           Identity
		multi              bool
		want               int
	}{
		{"org admin backup denied", "GET", "/api/v1/system/backups/backup/download", admin, true, 403},
		{"org admin mutation denied", "POST", "/api/v1/system/database/reset", admin, true, 403},
		{"org admin flags denied", "PATCH", "/api/v1/runtime-flags/auth", admin, true, 403},
		{"operator mutation allowed", "POST", "/api/v1/system/database/reset", operator, true, 204},
		{"single org admin allowed", "POST", "/api/v1/system/database/reset", admin, false, 204},
		{"member denied", "POST", "/api/v1/system/database/reset", member, false, 403},
		{"local flag enable allowed", "PATCH", "/api/v1/runtime-flags/auth", synthetic, false, 204},
		{"local actual auth flag enable allowed", "PATCH", "/api/v1/runtime-flags/features.auth", synthetic, false, 204},
		{"local bootstrap restart allowed", "POST", "/api/v1/system/restart", synthetic, false, 204},
		{"boot replay settings denied", "PATCH", "/api/v1/agent-profiles/profile", synthetic, false, 401},
		{"boot replay agent creation denied", "POST", "/api/v1/agents", synthetic, false, 401},
		{"PAT settings denied", "PATCH", "/api/v1/agent-profiles/profile", pat, false, 403},
		{"agent executor read denied", "GET", "/api/v1/executors/executor", pat, false, 403},
		{"local executor edit denied", "PATCH", "/api/v1/executors/executor/profiles/profile", synthetic, false, 401},
		{"local shared prompt edit denied", "PATCH", "/api/v1/prompts/prompt", synthetic, false, 401},
		{"agent workspace MCP edit denied", "POST", "/api/v1/workspaces/workspace/mcp-servers", pat, false, 403},
		{"agent MCP selection edit denied", "PUT", "/api/v1/tasks/task/mcp-selections", pat, false, 403},
		{"agent marketplace install denied", "POST", "/api/v1/workspaces/workspace/mcp-marketplace/install", pat, false, 403},
		{"human workspace MCP edit allowed", "POST", "/api/v1/workspaces/workspace/mcp-servers", admin, false, 204},
		{"PAT cannot mint admin account", "POST", "/api/v1/users", pat, false, 403},
		{"PAT cannot invite admin", "POST", "/api/v1/auth/invites", pat, false, 403},
		{"PAT cannot change password", "PATCH", "/api/v1/auth/password", pat, false, 403},
		{"PAT cannot download credentials", "GET", "/api/v1/system/backups/backup/download", pat, false, 403},
		{"member browser password allowed", "PATCH", "/api/v1/auth/password", member, false, 204},
		{"public invitation inspection retained", "GET", "/api/v1/auth/invites/preview", Identity{}, false, 204},
		{"human settings allowed", "PATCH", "/api/v1/agent-profiles/profile", admin, false, 204},
		{"foreign global settings denied", "PATCH", "/api/v1/agent-profiles/profile", admin, true, 403},
		{"ordinary task untouched", "GET", "/api/v1/tasks/123", member, true, 204},
		{"preview unavailable", "GET", "/port-proxy/session/3000/index.html", admin, false, 503},
		{"preview websocket unavailable", "GET", "/vscode/session/socket", admin, false, 503},
		{"canvas unavailable", "GET", "/api/v1/plugins/web-apps/runtime/token/index.html", admin, false, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) { SetOnGin(c, tc.identity) })
			r.Use(ControlPlaneBoundary(tc.multi))
			called := false
			r.Any(tc.path, func(c *gin.Context) { called = true; c.Status(204) })
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("X-Kandev-Interim-Settings-Interlock", "public-boot-token")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			if resp.Code != tc.want || called != (tc.want == 204) {
				t.Fatalf("status = %d, side effect = %v, want %d", resp.Code, called, tc.want)
			}
		})
	}
}

func TestGeneratedAPIHTMLCannotExecuteAtOperatorOrigin(t *testing.T) {
	for _, contentType := range []string{"text/html", "application/xhtml+xml", "image/svg+xml"} {
		t.Run(contentType, func(t *testing.T) {
			r := gin.New()
			r.Use(ControlPlaneBoundary(false))
			r.GET("/api/v1/attachments/file", func(c *gin.Context) {
				c.Data(http.StatusOK, contentType, []byte("<script>fetch('/api/v1/agents')</script>"))
			})
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/attachments/file", nil))
			if w.Header().Get("Content-Security-Policy") != "sandbox; default-src 'none'; frame-ancestors 'none'" ||
				w.Header().Get("Content-Disposition") != "attachment" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatalf("active API content retained operator-origin authority: %v", w.Header())
			}
		})
	}
}
