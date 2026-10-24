package rbac

import (
	"github.com/casbin/casbin/v3"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProvideEnforcer initializes the Casbin enforcer.
func ProvideEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	return NewEnforcer(db, "internal/config/rbac_model.conf")
}

var RBACSet = wire.NewSet(
	ProvideEnforcer,
)
