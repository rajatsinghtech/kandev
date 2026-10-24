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
		credentials := credentialSessionPath(p)
		if !settings && !install && !credentials {
			return
		}
		id, ok := FromGin(c)
		if ok && !multiTenancy && id.Synthetic && bootstrapAdministrationPath(p, c.Request.Method) {
			return
		}
		if status, reason := controlPlaneDenial(id, p, multiTenancy, settings, install, credentials); status != 0 {
			c.AbortWithStatusJSON(status, gin.H{"error": reason})
		}
	}
}

func controlPlaneDenial(id Identity, p string, multi, settings, install, credentials bool) (int, string) {
	if id.UserID == "" || id.Synthetic {
		return http.StatusUnauthorized, "authentication required"
	}
	if credentials && id.SessionID == "" {
		return http.StatusForbidden, "browser session required"
	}
	if !settings && !install {
		return 0, ""
	}
	if !HumanAdministrationSession(id) {
		return http.StatusForbidden, "human administration session required"
	}
	if multi && !accountAdministrationPath(p) && !id.Instance {
		return http.StatusForbidden, "installation authority required"
	}
	return 0, ""
}

func pathUnder(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, prefix+"/")
}

func unsafeContentPath(p string) bool {
	return pathUnder(p, "/port-proxy") || pathUnder(p, "/vscode") ||
		pathUnder(p, "/api/v1/plugins/web-apps/runtime")
}

func protectedSettingsPath(p string) bool {
	if accountAdministrationPath(p) {
		return true
	}
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

func accountAdministrationPath(p string) bool {
	return pathUnder(p, "/api/v1/users") || (pathUnder(p, "/api/v1/auth/invites") &&
		p != "/api/v1/auth/invites/preview" && p != "/api/v1/auth/invites/accept")
}

func credentialSessionPath(p string) bool {
	return pathUnder(p, "/api/v1/auth/password") || pathUnder(p, "/api/v1/auth/tokens") ||
		pathUnder(p, "/api/v1/auth/sessions")
}

func bootstrapAdministrationPath(p, method string) bool {
	return method == http.MethodPatch && (p == "/api/v1/runtime-flags/auth" || p == "/api/v1/runtime-flags/features.auth") ||
		method == http.MethodPost && p == "/api/v1/system/restart"
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
