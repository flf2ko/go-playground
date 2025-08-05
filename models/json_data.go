package models

import (
	"time"
)

type JSONRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	URL       string    `json:"url" gorm:"not null;size:2048"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type FetchRequest struct {
	Link string `json:"link" form:"link" binding:"required"`
}

type FetchResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *JSONRecord `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}