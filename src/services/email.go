package services

import (
    "bank-service/src/repositories"
    "fmt"
    "github.com/sirupsen/logrus"
    gomail "gopkg.in/gomail.v2"
)

type EmailService struct {
    dialer  *gomail.Dialer
    from    string
    logger  *logrus.Logger
}

func NewEmailService(host string, port int, user, pass, from string, logger *logrus.Logger) *EmailService {
    return &EmailService{
        dialer:  gomail.NewDialer(host, port, user, pass),
        from:    from,
        logger:  logger,
    }
}

func (s *EmailService) SendPaymentNotification(userRepo *repositories.UserRepository, userID uint, amount string) error {
    // Извлекаем пользователя из репозитория по ID
    user, err := userRepo.GetByID(userID)
    if err != nil {
        s.logger.Errorf("Failed to get user by ID: %v", err)
        return err
    }

    // Создаем новое сообщение
    m := gomail.NewMessage()
    m.SetHeader("From", s.from)            // Устанавливаем адрес отправителя
    m.SetHeader("To", user.Email)           // Устанавливаем адрес получателя
    m.SetHeader("Subject", "Платеж по кредиту") // Устанавливаем тему сообщения
    m.SetBody("text/html", fmt.Sprintf("Сумма платежа: %s", amount)) // Устанавливаем тело сообщения

    // Отправляем сообщение
    if err := s.dialer.DialAndSend(m); err != nil {
        s.logger.Errorf("Failed to send email: %v", err)
        return err
    }

    return nil // Успешная отправка
}

