// internal/middleware/rbac.go
package auth

import (
	"net/http"
	"seo-backend/internal/infrastructure/db/repositories/rbac"
)

// Require adalah factory middleware
// module  → "draft" | "histories" | "product" | "user_management"
// action  → "view" | "create" | "edit" | "delete" | "publish" | "inject" | "assign_role"
// scope   → "team" | "global"
func Require(cache *rbac.PermissionCache, module, action, scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims := GetClaims(ctx)

			if claims == nil {
				writeError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			perms, err := cache.Get(ctx, claims.Role)
			if err != nil {
				writeError(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !rbac.HasPermission(perms, module, action, scope) {
				writeError(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ---- Shortcut per resource ----

func DraftView(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "view", "team")
}
func DraftViewGlobal(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "view", "global")
}
func DraftCreate(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "create", "team")
}
func DraftEdit(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "edit", "team")
}
func DraftDelete(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "delete", "team")
}
func DraftPublish(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "draft", "publish", "team")
}

func ProductView(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "view", "team")
}
func ProductViewGlobal(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "view", "global")
}
func ProductCreate(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "create", "team")
}
func ProductEdit(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "edit", "team")
}
func ProductDelete(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "delete", "team")
}
func ProductInject(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "product", "inject", "team")
}

func UserManage(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "user_management", "manage", "team")
}
func UserManageGlobal(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "user_management", "manage", "global")
}
func UserAssignRole(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "user_management", "assign_role", "team")
}

// History — sesuai schema: view(team), view(global), delete(global)
func HistoryView(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "histories", "view", "team")
}
func HistoryViewGlobal(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "histories", "view", "global")
}
func HistoryDeleteGlobal(c *rbac.PermissionCache) func(http.Handler) http.Handler {
	return Require(c, "histories", "delete", "global")
}
