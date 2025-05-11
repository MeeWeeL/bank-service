package models

import "time"

type Transaction struct {
    ID            uint      `json:"id"`
    FromAccountID uint      `json:"from_account_id,omitempty"`
    ToAccountID   uint      `json:"to_account_id"`
    Amount        float64   `json:"amount"`
    Currency      string    `json:"currency"`
    CreatedAt     time.Time `json:"created_at"`
}