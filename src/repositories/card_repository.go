package repositories

import (
	"bank-service/src/models"
	"database/sql"
	"errors"
	
	"github.com/sirupsen/logrus"
)

var (
	ErrCardNotFound = errors.New("card not found")
)

type CardRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewCardRepository(db *sql.DB, logger *logrus.Logger) *CardRepository {
	return &CardRepository{db: db, logger: logger}
}

func (r *CardRepository) Create(card *models.Card) error {
	query := `INSERT INTO cards 
		(user_id, account_id, encrypted_data, hmac, cvv_hash) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at`
		
	return r.db.QueryRow(query,
		card.UserID,
		card.AccountID,
		card.EncryptedData,
		card.Hmac,
		card.CvvHash,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
}

func (r *CardRepository) GetByIDAndUser(cardID, userID uint) (*models.Card, error) {
	card := &models.Card{}
	query := `SELECT id, user_id, account_id, encrypted_data, hmac, created_at, updated_at 
		FROM cards WHERE id = $1 AND user_id = $2`
	
	err := r.db.QueryRow(query, cardID, userID).Scan(
		&card.ID,
		&card.UserID,
		&card.AccountID,
		&card.EncryptedData,
		&card.Hmac,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCardNotFound
	}
	return card, err
}

func (r *CardRepository) GetByUser(userID uint) ([]models.Card, error) {
	query := `SELECT id, account_id, encrypted_data, hmac, created_at, updated_at 
		FROM cards WHERE user_id = $1`
		
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(
			&card.ID,
			&card.AccountID,
			&card.EncryptedData,
			&card.Hmac,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}