package rules

import (
	"go-starterkit/internal/modules/auth"
	user "go-starterkit/internal/modules/users"

	"gorm.io/gen"
)

// RegisterAllRules collecting all module relationship rules for GORM Gen
func RegisterAllRules(g *gen.Generator) {
	rules := []func(*gen.Generator){
		auth.RegisterRulesAuth,
		user.RegisterRules,
	}

	for _, rule := range rules {
		rule(g)
	}
}
