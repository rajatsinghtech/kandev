package authn

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// HumanAdministrationSession excludes agent identities and programmatic tokens.
func HumanAdministrationSession(id Identity) bool {
	return !id.Synthetic && id.UserID != "" && id.SessionID != "" &&
		(id.IsAdmin() || id.Role == "owner" || id.Instance)
}

// ControlPlaneBoundary runs after authentication and before route handlers.
// Generated applications remain unavailable until they have a separate site.
func ControlPlaneBoundary(multiTenancy bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := path.Clean(c.Request.URL.Path)
		if pathUnder(p, "/api/v1") {
			c.Writer = &apiContentWriter{ResponseWriter: c.Writer}
		}
		if unsafeContentPath(p) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "active content requires an isolated origin", "error_code": "isolated_origin_required",
			})
			return
		}
		settings := protectedSettingsPath(p)
		if c.Request.Method == http.MethodGet && p == "/api/v1/agents" {
			settings = false
		}
		install := installationPath(p, c.Request.Method)
		if !settings && !install {
			return
		}
		id, ok := FromGin(c)
		if !ok || (settings && id.Synthetic) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		if settings && !HumanAdministrationSession(id) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "human administration session required"})
			return
		}
		if multiTenancy && (id.Synthetic || !id.Instance) || !multiTenancy && !id.IsAdmin() && id.Role != "owner" && !id.Instance {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "installation authority required"})
		}
	}
}

func pathUnder(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, prefix+"/")
}

func unsafeContentPath(p string) bool {
	return pathUnder(p, "/port-proxy") || pathUnder(p, "/vscode") ||
		pathUnder(p, "/api/v1/plugins/web-apps/runtime")
}

func protectedSettingsPath(p string) bool {
	for _, prefix := range []string{"agents", "agent-profiles", "agent-models", "agent-command-preview", "agent-install", "agent-update", "executors", "executor-profiles", "editors", "prompts"} {
		if pathUnder(p, "/api/v1/"+prefix) {
			return true
		}
	}
	if pathUnder(p, "/api/v1") {
		for _, segment := range strings.Split(p, "/") {
			if segment == "mcp-servers" || segment == "mcp-selections" || segment == "mcp-marketplace" {
				return true
			}
		}
	}
	return false
}

func installationPath(p, method string) bool {
	write := method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
	if pathUnder(p, "/api/v1/system") {
		return write || pathUnder(p, "/api/v1/system/backups") || pathUnder(p, "/api/v1/system/logs")
	}
	for _, prefix := range []string{"runtime-flags", "executors", "executor-profiles", "environments", "editors", "prompts"} {
		if write && pathUnder(p, "/api/v1/"+prefix) {
			return true
		}
	}
	return false
}
