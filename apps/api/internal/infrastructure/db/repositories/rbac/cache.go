// internal/infrastructure/rbac/cache.go
package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"seo-backend/internal/infrastructure/db/repositories"
	"time"

	"github.com/go-redis/redis/v8"
)

type PermissionCache struct {
	redis *redis.Client
	cache map[string][]repositories.Permission // key: role name
	repo  *repositories.RBACRepository
	ttl   time.Duration
}

func NewPermissionCache(redis *redis.Client, repo *repositories.RBACRepository, ttl time.Duration) *PermissionCache {
	return &PermissionCache{
		redis: redis,
		repo:  repo,
		ttl:   ttl,
	}
}

const permissionCacheKey = "permissions:role:%s"

// Get → ambil dari cache, kalau belum ada query DB
func (c *PermissionCache) Get(ctx context.Context, roleName string) ([]repositories.Permission, error) {
	key := fmt.Sprintf(permissionCacheKey, roleName)

	// 1. Try Redis
	val, err := c.redis.Get(ctx, key).Bytes()
	if err == nil {
		var perms []repositories.Permission
		if err := json.Unmarshal(val, &perms); err != nil {
			return nil, fmt.Errorf("unmarshal cache: %w", err)
		}
		return perms, nil
	}

	if !errors.Is(err, redis.Nil) {
		// Redis error (bukan cache miss) — log tapi jangan block
		log.Printf("[PermissionCache] redis get error: %v", err)
	}

	// 2. Cache miss → query DB
	perms, err := c.repo.GetPermissionsByRole(ctx, roleName)
	if err != nil {
		return nil, err
	}

	// 3. Store to Redis
	data, err := json.Marshal(perms)
	if err != nil {
		return nil, fmt.Errorf("marshal cache: %w", err)
	}

	if err := c.redis.Set(ctx, key, data, c.ttl).Err(); err != nil {
		// Jangan return error, data sudah ada dari DB
		log.Printf("[PermissionCache] redis set error: %v", err)
	}

	return perms, nil
}

// Invalidate → hapus cache role tertentu (misal setelah assign role)
func (c *PermissionCache) Invalidate(ctx context.Context, roleName string) error {
	key := fmt.Sprintf(permissionCacheKey, roleName)
	return c.redis.Del(ctx, key).Err()
}

// HasPermission → cek apakah slice permission punya module+action+scope tertentu
func HasPermission(perms []repositories.Permission, module, action, scope string) bool {
	for _, p := range perms {
		if p.Module == module && p.Action == action && p.Scope == scope {
			return true
		}
	}
	return false
}
