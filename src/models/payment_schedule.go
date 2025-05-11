package models

import "time"

type PaymentSchedule struct {
    ID        uint      `json:"id"`
    CreditID  uint      `json:"credit_id"`
    DueDate   time.Time `json:"due_date"`
    Amount    float64   `json:"amount"`
    Paid      bool      `json:"paid"`
    CreatedAt time.Time `json:"created_at"`
}
