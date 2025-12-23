package models

import (
	"time"
)

type TodoItem struct {
	ID          uint      `gorm:"primarykey"`
	Title       string    `gorm:"size:50;unique;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	IsDone      bool      `gorm:"default:false" json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
