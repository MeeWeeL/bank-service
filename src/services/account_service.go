package services

import (
    "bank-service/src/models"
    "bank-service/src/repositories"
    "errors"
    "github.com/sirupsen/logrus"
)

type AccountService struct {
    accountRepo     *repositories.AccountRepository
    transactionRepo *repositories.TransactionRepository
    logger          *logrus.Logger
}

func NewAccountService(
    accountRepo *repositories.AccountRepository,
    transactionRepo *repositories.TransactionRepository,
    logger *logrus.Logger,
) *AccountService {
    return &AccountService{
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
        logger:          logger,
    }
}

func (s *AccountService) CreateAccount(userID uint) (*models.Account, error) {
    account := &models.Account{
        UserID:   userID,
        Currency: "RUB",
    }

    if err := s.accountRepo.Create(account); err != nil {
        return nil, err
    }
    
    return account, nil
}

func (s *AccountService) GetAccount(accountID uint) (*models.Account, error) {
    return s.accountRepo.GetByID(accountID)
}

func (s *AccountService) Transfer(fromAccountID, toAccountID uint, amount float64) error {
    tx, err := s.accountRepo.BeginTx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    if err := s.accountRepo.UpdateBalanceTx(tx, fromAccountID, -amount); err != nil {
        return err
    }

    if err := s.accountRepo.UpdateBalanceTx(tx, toAccountID, amount); err != nil {
        return err
    }

    if _, err := tx.Exec(
        `INSERT INTO transactions (from_account_id, to_account_id, amount, currency)
         VALUES ($1, $2, $3, 'RUB')`,
        fromAccountID, toAccountID, amount,
    ); err != nil {
        return err
    }

    return tx.Commit()
}


func (s *AccountService) Deposit(accountID uint, amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    return s.accountRepo.UpdateBalance(accountID, amount)
}

func (s *AccountService) GetByIDAndUser(accountID, userID uint) (*models.Account, error) {
    return s.accountRepo.GetByIDAndUser(accountID, userID)
}