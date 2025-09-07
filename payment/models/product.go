package models

import (
	"time"
)

type Product struct {
	ProductID string `json:"productId"`
	Price     int64  `json:"price"`
	Currency  string `json:"currency"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;"`
}

type ProductInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
