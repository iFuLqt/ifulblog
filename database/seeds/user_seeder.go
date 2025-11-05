package seeds

import (
	"ifulblog/internal/core/domain/model"
	"ifulblog/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating password hash")
	}

	admin := model.User{
		Name:     "Admin",
		Email:    "Admin@example.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seeding role admin")
	} else {
		log.Info().Msg("Admin role seeding succesfully")
	}

}
