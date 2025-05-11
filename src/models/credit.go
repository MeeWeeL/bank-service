package models

import "time"

type Credit struct {
    ID         uint      `json:"id"`
    UserID     uint      `json:"user_id"`
    AccountID  uint      `json:"account_id"`
    Amount     float64   `json:"amount"`
    Rate       float64   `json:"rate"`
    Period     int       `json:"period"` // months
    CreatedAt  time.Time `json:"created_at"`
    Status     string    `json:"status"`
}