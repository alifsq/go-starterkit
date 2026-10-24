package rbac

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, fmt.Errorf("rbac : database connection nil")
	}
	if modelPath == "" {
		return nil, fmt.Errorf("rbac : model path is empty")
	}
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("rbac : create GORM Adapter : %w", err)
	}
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("rbac : create enforcer : %w", err)

	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	// Ensure default permissions exist even when the database already has policies.
	for role, permissions := range RolePermission {
		for _, perm := range permissions {
			hasPolicy, err := enforcer.HasPolicy(role, perm.Path, perm.Method)
			if err != nil {
				return nil, fmt.Errorf("rbac: check default policy for role %q: %w", role, err)
			}
			if hasPolicy {
				continue
			}
			if _, err := enforcer.AddPolicy(role, perm.Path, perm.Method); err != nil {
				return nil, fmt.Errorf("rbac: add default policy for role %q: %w", role, err)
			}
		}
	}
	return enforcer, nil
}
