package config

import (
	"github.com/Alieksieiev0/feed-service/internal/database"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	err = database.Setup(db)
	if err != nil {
		return nil, err
	}

	return db, err
}
