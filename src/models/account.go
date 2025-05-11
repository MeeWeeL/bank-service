package models

import "time"

type Account struct {
    ID        uint      `json:"id"`
    UserID    uint      `json:"user_id" validate:"required"`
    Balance   float64   `json:"balance" validate:"gte=0"`
    Currency  string    `json:"currency" validate:"required,eq=RUB"`
    CreatedAt time.Time `json:"created_at"`
}