package repositories

import (
	"bank-service/src/models"
	"database/sql"
	"errors"
	
	"github.com/sirupsen/logrus"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewAccountRepository(db *sql.DB, logger *logrus.Logger) *AccountRepository {
	return &AccountRepository{db: db, logger: logger}
}

func (r *AccountRepository) Create(account *models.Account) error {
	return r.db.QueryRow(
		`INSERT INTO accounts (user_id, currency) 
		 VALUES ($1, $2) RETURNING id, created_at`,
		account.UserID, account.Currency,
	).Scan(&account.ID, &account.CreatedAt)
}

func (r *AccountRepository) GetByIDAndUser(accountID, userID uint) (*models.Account, error) {
	account := &models.Account{}
	err := r.db.QueryRow(
		`SELECT id, user_id, balance, currency, created_at 
		 FROM accounts WHERE id = $1 AND user_id = $2`,
		accountID, userID,
	).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt)
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	return account, err
}

func (r *AccountRepository) GetByID(accountID uint) (*models.Account, error) {
	account := &models.Account{}
	err := r.db.QueryRow(
		`SELECT id, user_id, balance, currency, created_at 
		 FROM accounts WHERE id = $1`,
		accountID,
	).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt)
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	return account, err
}

func (r *AccountRepository) UpdateBalance(accountID uint, amount float64) error {
    _, err := r.db.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount,
        accountID,
    )
    return err
}

func (r *AccountRepository) UpdateBalanceTx(tx *sql.Tx, accountID uint, amount float64) error {
    res, err := tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount,
        accountID,
    )
    if err != nil {
        return err
    }
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return ErrAccountNotFound
    }
    return nil
}

func (r *AccountRepository) BeginTx() (*sql.Tx, error) {
    return r.db.Begin()
}

func (r *CreditRepository) GetByAccountID(accountID uint) ([]models.Credit, error) {
	query := `SELECT id, user_id, account_id, amount, rate, period, created_at, status 
	          FROM credits WHERE account_id = $1`
	rows, err := r.db.Query(query, accountID)
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
