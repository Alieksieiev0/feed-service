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
	ID        string         `gorm:"type:uuid" json:"id"`
	CreatedAt time.Time      `                 json:"created_at"`
	UpdatedAt time.Time      `                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"     json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

type User struct {
	Base
	Username    string `gorm:"default:null;not null;unique;"`
	Password    string `gorm:"default:null;not null;"`
	Email       string `gorm:"default:null;not null;unique;"`
	Subscribers []User `gorm:"many2many:user_subscribers"`
	Posts       []Post
}

type Post struct {
	Base
	Title  string `json:"title"   gorm:"not null; default:null;"`
	Body   string `json:"body"    gorm:"not null; default:null;"`
	UserId string `json:"user_id" gorm:"type:uuid; not null; default:null;"`
}
