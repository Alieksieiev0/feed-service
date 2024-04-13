package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity interface {
	BeforeCreate(tx *gorm.DB) error
}

type Base struct {
	ID        string `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

type User struct {
	Base
	Username   string `gorm:"default:null;not null;unique;"`
	Password   string `gorm:"default:null;not null;"`
	Email      string `gorm:"default:null;not null;unique;"`
	Subcribers []User `gorm:"many2many:user_subcribers"`
	Posts      []Post
}

type Post struct {
	Base
	Title  string `json:"title"   gorm:"not null; default:null;"`
	Body   string `json:"body"    gorm:"not null; default:null;"`
	UserId string `json:"user_id" gorm:"not null; default:null;"`
}
