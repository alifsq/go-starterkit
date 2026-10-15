package config

import (
	"go-starterkit/internal/database/rules"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func RunGenerator(db *gorm.DB, outPath string, modelPath string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       outPath,
		ModelPkgPath:  modelPath,
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	// RULE GLOBAL: Otomatis pasang tag `<-:create` untuk SEMUA kolom created_at di tabel apa pun
	g.WithOpts(
		gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
			return tag.Append("<-:create").Append("autoCreateTime")
		}),
		gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
			return tag.Append("<-:update").Append("autoUpdateTime")
		}),
	)
	
	// All rules
	rules.RegisterAllRules(g)
	g.GenerateModel("schema_migrations")

	g.Execute()
}
