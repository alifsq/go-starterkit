package auth

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func RegisterRules(g *gen.Generator) {
	// Initialize
	userModel := g.GenerateModel("users")

	// Register Relation In Here !!!
	refreshTokenWithRel := g.GenerateModel("refresh_tokens",
		gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{
			RelatePointer: true,
			GORMTag:       field.GormTag{"foreignKey": []string{"UserID"}, "references": []string{"ID"}},
		}),
	)

	g.ApplyBasic(refreshTokenWithRel)
}
