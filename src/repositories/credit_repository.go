package repositories

import (
	"bank-service/src/models"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	ErrCreditNotFound = errors.New("credit not found")
)

type CreditRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewCreditRepository(db *sql.DB, logger *logrus.Logger) *CreditRepository {
	return &CreditRepository{db: db, logger: logger}
}

func (r *CreditRepository) Create(credit *models.Credit) error {
	query := `INSERT INTO credits (user_id, account_id, amount, rate, period, status) 
          VALUES ($1, $2, $3, $4, $5, $6) 
          RETURNING id, created_at`
	r.logger.Infof("Executing query: %s with values: userID=%d, accountID=%d, amount=%f, rate=%f, period=%d, status=%s",
		query, credit.UserID, credit.AccountID, credit.Amount, credit.Rate, credit.Period, credit.Status)
	return r.db.QueryRow(query,
		credit.UserID,
		credit.AccountID,
		credit.Amount,
		credit.Rate,
		credit.Period,
		credit.Status,
	).Scan(&credit.ID, &credit.CreatedAt)
}

func (r *CreditRepository) GetByIDAndUser(creditID, userID uint) (*models.Credit, error) {
	credit := &models.Credit{}
	query := `SELECT id, user_id, account_id, amount, rate, period, created_at, status 
          FROM credits WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, creditID, userID).Scan(
		&credit.ID,
		&credit.UserID,
		&credit.AccountID,
		&credit.Amount,
		&credit.Rate,
		&credit.Period,
		&credit.CreatedAt,
		&credit.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCreditNotFound
	}
	return credit, err
}

func (r *CreditRepository) UpdateStatus(creditID uint, status string) error {
    _, err := r.db.Exec("UPDATE credits SET status = $1 WHERE id = $2", status, creditID)
    return err
}

func (r *CreditRepository) GetByUserID(userID uint) ([]models.Credit, error) {
	query := `SELECT id, user_id, account_id, amount, rate, period, created_at, status 
	          FROM credits WHERE user_id = \$1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []models.Credit
	for rows.Next() {
		var c models.Credit
		if err := rows.Scan(&c.ID, &c.UserID, &c.AccountID, &c.Amount, &c.Rate, &c.Period, &c.CreatedAt, &c.Status); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}
	return credits, nil
}
