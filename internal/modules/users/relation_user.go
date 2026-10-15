package user

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func RegisterRules(g *gen.Generator) {

	// ==== Initialize ====//
	refreshTokenModel := g.GenerateModel("refresh_tokens")

	// ==== Register Relation In Here !!! === //

	// HasMany -> To refresh_tokens
	userModelWithRel := g.GenerateModel("users",
		gen.FieldRelate(field.HasMany, "RefreshTokens", refreshTokenModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"UserID"}, "references": []string{"ID"}},
		}),
	)
	g.ApplyBasic(userModelWithRel)
}
