package model

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	ID        uuid.UUID `gorm:"column:ID;type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	CreatedBy string `gorm:"column:created_by"`
	Image     []byte `gorm:"size:1048576"` // 1MB
}
