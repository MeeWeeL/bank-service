package services

import (
	"bank-service/src/models"
	"bank-service/src/repositories"
	"errors"
	"math"
	"time"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CreditService struct {
	creditRepo          *repositories.CreditRepository
	paymentScheduleRepo *repositories.PaymentScheduleRepository
	accountService      *AccountService
	logger              *logrus.Logger
    cbrService          *CBRService
	userRepo            *repositories.UserRepository
	emailService        *EmailService       
}

func NewCreditService(
	creditRepo *repositories.CreditRepository,
	paymentScheduleRepo *repositories.PaymentScheduleRepository,
	accountService *AccountService,
    cbrService *CBRService,
	logger *logrus.Logger,
	userRepo            *repositories.UserRepository,
	emailService        *EmailService,
) *CreditService {
	return &CreditService{
		creditRepo:          creditRepo,
		paymentScheduleRepo: paymentScheduleRepo,
		accountService:      accountService,
        cbrService:       cbrService,
		logger:              logger,
		userRepo: userRepo,
		emailService: emailService,
	}
}

// Оформление кредита с расчетом аннуитетных платежей
func (s *CreditService) CreateCredit(userID, accountID uint, amount float64, rate float64, period int) (*models.Credit, error) {
    if rate <= 0 {
        keyRate, err := s.cbrService.GetKeyRate()
        if err != nil {
            s.logger.Warnf("Using default rate 10%%, failed to get CBR rate: %v", err)
            rate = 10.0 // дефолтная ставка при ошибке
        } else {
            rate = keyRate
        }
    }

	if amount <= 0 || rate <= 0 || period <= 0 {
		return nil, errors.New("invalid credit parameters")
	}

	credit := &models.Credit{
		UserID:    userID,
		AccountID: accountID,
		Amount:    amount,
		Rate:      rate,
		Period:    period,
		Status:    "active",
	}

	err := s.creditRepo.Create(credit)
	if err != nil {
		return nil, err
	}

	// Генерация графика платежей
	err = s.generatePaymentSchedule(credit)
	if err != nil {
		return nil, err
	}

	return credit, nil
}

// Расчет аннуитетного платежа и создание графика платежей
func (s *CreditService) generatePaymentSchedule(credit *models.Credit) error {
	// Ежемесячная ставка в долях (rate - годовая в процентах)
	monthlyRate := credit.Rate / 100 / 12
	P := credit.Amount
	n := float64(credit.Period)

	// Аннуитетный платеж A = P * (r * (1+r)^n) / ((1+r)^n - 1)
	A := P * (monthlyRate * math.Pow(1+monthlyRate, n)) / (math.Pow(1+monthlyRate, n) - 1)
	A = math.Round(A*100) / 100 // округление до копеек

	for i := 1; i <= credit.Period; i++ {
		dueDate := time.Now().AddDate(0, i, 0)
		schedule := &models.PaymentSchedule{
			CreditID: credit.ID,
			DueDate:  dueDate,
			Amount:   A,
			Paid:     false,
		}
		err := s.paymentScheduleRepo.Create(schedule)
		if err != nil {
			return err
		}
	}

	return nil
}

// Получение графика платежей по кредиту
func (s *CreditService) GetPaymentSchedule(userID, creditID uint) ([]models.PaymentSchedule, error) {
	credit, err := s.creditRepo.GetByIDAndUser(creditID, userID)
	if err != nil {
		return nil, err
	}

	return s.paymentScheduleRepo.GetByCreditID(credit.ID)
}

// Обработка просроченных платежей (начисление штрафов, попытка списания)
func (s *CreditService) ProcessOverduePayments() error {
	now := time.Now()
	overdueSchedules, err := s.paymentScheduleRepo.GetOverdueUnpaidSchedules(now)
	if err != nil {
		return err
	}

	for _, schedule := range overdueSchedules {
		credit, err := s.creditRepo.GetByIDAndUser(schedule.CreditID, 0) // 0 - без проверки пользователя
		if err != nil {
			s.logger.WithError(err).Warnf("Credit not found for schedule %d", schedule.ID)
			continue
		}

		account, err := s.accountService.GetAccount(credit.AccountID)
		if err != nil {
			s.logger.WithError(err).Warnf("Account not found for credit %d", credit.ID)
			continue
		}

		// Штраф +10%
		amountWithPenalty := schedule.Amount * 1.10

		if account.Balance >= amountWithPenalty {
			// Списываем средства
			err = s.accountService.Transfer(credit.AccountID, 0, amountWithPenalty) // 0 - внешний счет (например банк)
			if err != nil {
				s.logger.WithError(err).Warnf("Failed to transfer payment for schedule %d", schedule.ID)
				continue
			}

			// Отмечаем платеж как оплаченный
			err = s.paymentScheduleRepo.MarkAsPaid(schedule.ID)
			if err != nil {
				s.logger.WithError(err).Warnf("Failed to mark schedule %d as paid", schedule.ID)
				continue
			}
		} else {
			s.logger.Infof("Insufficient funds for credit %d, schedule %d - penalty applied", credit.ID, schedule.ID)
			if err == nil {
				err = s.emailService.SendPaymentNotification(
					s.userRepo,          // Передаем userRepo
					credit.UserID,      // Передаем идентификатор пользователя
					fmt.Sprintf("%.2f RUB", amountWithPenalty), // Передаем сумму
				)
				if err != nil {
					// Обработка ошибки при отправке уведомления
				}
			}
		}
	}

	return nil
}

func (s *CreditService) GetCreditsByAccount(accountID uint) ([]models.Credit, error) {
	return s.creditRepo.GetByAccountID(accountID)
}

func (s *CreditService) AccountService() *AccountService {
    return s.accountService
}
