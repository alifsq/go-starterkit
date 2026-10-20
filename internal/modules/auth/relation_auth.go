package auth

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

func RegisterRulesAuth(g *gen.Generator) {
	// Initialize
	userModel := g.GenerateModel("users")
	emailTokenModel := g.GenerateModel("email_verification_tokens")
	passwordTokenModel := g.GenerateModel("password_reset_tokens")

	// Register Relation In Here !!!
	refreshTokenWithRel := g.GenerateModel("refresh_tokens",
		gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{
			RelatePointer: true,
			GORMTag:       field.GormTag{"foreignKey": []string{"UserID"}, "references": []string{"ID"}},
		}),
	)

	emailTokenWithRel := g.GenerateModel("email_verification_tokens",
		gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{
			RelatePointer: true,
			GORMTag:       field.GormTag{"foreignKey": []string{"UserID"}, "references": []string{"ID"}},
		}),
	)

	passwordTokenWithRel := g.GenerateModel("password_reset_tokens",
		gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"UserID"}, "references": []string{"ID"}},
		}),
	)

	// User
	userWithRel := g.GenerateModel("users",
		gen.FieldRelate(field.HasMany, "EmailVerificationTokens", emailTokenModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasMany, "PasswordResetTokens", passwordTokenModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
	)

	g.ApplyBasic(refreshTokenWithRel, emailTokenWithRel, passwordTokenWithRel, userWithRel)
}
