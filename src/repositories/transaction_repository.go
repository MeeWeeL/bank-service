package repositories

import (
    "bank-service/src/models"
    "database/sql"
    "time"
    "github.com/sirupsen/logrus"
)

type TransactionRepository struct {
    db     *sql.DB
    logger *logrus.Logger
}

func NewTransactionRepository(db *sql.DB, logger *logrus.Logger) *TransactionRepository {
    return &TransactionRepository{db: db, logger: logger}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
    return r.db.QueryRow(
        `INSERT INTO transactions (from_account_id, to_account_id, amount, currency)
         VALUES ($1, $2, $3, $4)
         RETURNING id, created_at`,
        transaction.FromAccountID,
        transaction.ToAccountID,
        transaction.Amount,
        transaction.Currency,
    ).Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *TransactionRepository) SumIncome(userID uint, start, end time.Time) (float64, error) {
    var income float64
    err := r.db.QueryRow(
        `SELECT COALESCE(SUM(amount), 0) 
         FROM transactions t
         JOIN accounts a ON t.to_account_id = a.id
         WHERE a.user_id = \$1 AND t.created_at BETWEEN \$2 AND \$3`,
        userID, start, end,
    ).Scan(&income)
    return income, err
}

func (r *TransactionRepository) SumExpenses(userID uint, start, end time.Time) (float64, error) {
    var expenses float64
    err := r.db.QueryRow(
        `SELECT COALESCE(SUM(amount), 0)
         FROM transactions t
         JOIN accounts a ON t.from_account_id = a.id
         WHERE a.user_id = \$1 AND t.created_at BETWEEN \$2 AND \$3`,
        userID, start, end,
    ).Scan(&expenses)
    return expenses, err
}