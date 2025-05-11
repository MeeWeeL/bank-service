package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type Card struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id" validate:"required"`
	AccountID     uint      `json:"account_id" validate:"required"`
	EncryptedData string    `json:"encrypted_data"` // PGP encrypted (number + expiry)
	Hmac          string    `json:"hmac"`           // HMAC-SHA256 of encrypted data
	CvvHash       string    `json:"-"`              // bcrypt hash
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *Card) ValidateLuhn(number string) bool {
	sum := 0
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}

func (c *Card) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}