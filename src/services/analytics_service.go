package services

import (
	"bank-service/src/repositories"
	"time"
)

type AnalyticsService struct {
	transactionRepo *repositories.TransactionRepository
	creditRepo      *repositories.CreditRepository
	accountRepo     *repositories.AccountRepository
	paymentScheduleRepo *repositories.PaymentScheduleRepository 
}

func NewAnalyticsService(
	transactionRepo *repositories.TransactionRepository,
	creditRepo *repositories.CreditRepository,
	accountRepo *repositories.AccountRepository,
	paymentScheduleRepo *repositories.PaymentScheduleRepository,
) *AnalyticsService {
	return &AnalyticsService{
		transactionRepo: transactionRepo,
		creditRepo:      creditRepo,
		accountRepo:     accountRepo,
		paymentScheduleRepo: paymentScheduleRepo,
	}
}

func (s *AnalyticsService) GetMonthlyIncomeExpenses(userID uint, year int, month time.Month) (income, expenses float64, err error) {
    start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
    end := start.AddDate(0, 1, 0)

    income, err = s.transactionRepo.SumIncome(userID, start, end)
    if err != nil {
        return 0, 0, err
    }

    expenses, err = s.transactionRepo.SumExpenses(userID, start, end)
    if err != nil {
        return 0, 0, err
    }

    return income, expenses, nil
}


// Аналитика кредитной нагрузки
func (s *AnalyticsService) GetCreditLoad(userID uint) (float64, error) {
	credits, err := s.creditRepo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, c := range credits {
		if c.Status == "active" {
			total += c.Amount
		}
	}
	return total, nil
}


// Прогноз баланса на N дней (учет запланированных платежей)
func (s *AnalyticsService) PredictBalance(accountID uint, days int) (float64, error) {
    // Получаем текущий баланс аккаунта
    account, err := s.accountRepo.GetByID(accountID)
    if err != nil {
        return 0, err
    }

    now := time.Now()
    endDate := now.AddDate(0, 0, days)

    // Получаем кредиты по аккаунту
    credits, err := s.creditRepo.GetByAccountID(accountID)
    if err != nil {
        return 0, err
    }

    // Суммируем все запланированные платежи по кредитам за период
    var totalPayments float64 = 0
    for _, credit := range credits {
        // Получаем график платежей по кредиту
        schedules, err := s.paymentScheduleRepo.GetByCreditID(credit.ID)
        if err != nil {
            return 0, err
        }

        for _, sched := range schedules {
            if sched.DueDate.After(now) && !sched.DueDate.After(endDate) && !sched.Paid {
                totalPayments += sched.Amount
            }
        }
    }

    // Прогнозируем баланс: текущий баланс - ожидаемые платежи
    predictedBalance := account.Balance - totalPayments

    return predictedBalance, nil
}

