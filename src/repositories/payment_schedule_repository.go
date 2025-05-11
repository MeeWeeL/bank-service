package repositories

import (
	"bank-service/src/models"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type PaymentScheduleRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPaymentScheduleRepository(db *sql.DB, logger *logrus.Logger) *PaymentScheduleRepository {
	return &PaymentScheduleRepository{db: db, logger: logger}
}

func (r *PaymentScheduleRepository) Create(schedule *models.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (credit_id, due_date, amount, paid) 
          VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(query,
		schedule.CreditID,
		schedule.DueDate,
		schedule.Amount,
		schedule.Paid,
	).Scan(&schedule.ID, &schedule.CreatedAt)
}

func (r *PaymentScheduleRepository) GetByCreditID(creditID uint) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, due_date, amount, paid, created_at 
          FROM payment_schedules WHERE credit_id = $1 ORDER BY due_date`
	rows, err := r.db.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.DueDate, &s.Amount, &s.Paid, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *PaymentScheduleRepository) MarkAsPaid(scheduleID uint) error {
	_, err := r.db.Exec("UPDATE payment_schedules SET paid = TRUE WHERE id = $1", scheduleID)
	return err
}

func (r *PaymentScheduleRepository) GetOverdueUnpaidSchedules(before time.Time) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, due_date, amount, paid, created_at 
          FROM payment_schedules WHERE due_date < $1 AND paid = FALSE`
	rows, err := r.db.Query(query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.DueDate, &s.Amount, &s.Paid, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}
