package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint64         `gorm:"primaryKey; auto_increment; index;" json:"taskId"`
	Title       string         `gorm:"not null"`
	Description *string        `gorm:"default:null"`
	Status      bool           `gorm:"default:false; not null"`
	CreatedAt   *time.Time     `gorm:"default:now()"  json:"createdAt"`
	UpdatedAt   *time.Time     `gorm:"default:now()"  json:"updateAt"`
	DeleteAt    gorm.DeletedAt `gorm:"default:null" json:"-"`
}
