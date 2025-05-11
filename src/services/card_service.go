package services

import (
    "bank-service/src/models"
    "bank-service/src/repositories"
	"bank-service/src/crypto"
	hmacpkg "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math/rand"
    "time"
	
    "golang.org/x/crypto/bcrypt"
    "github.com/ProtonMail/go-crypto/openpgp"
    "github.com/sirupsen/logrus"
)

type CardService struct {
	cardRepo    *repositories.CardRepository
	accountRepo *repositories.AccountRepository
	pgpEntity   *openpgp.Entity
	logger      *logrus.Logger
}

func NewCardService(
	cardRepo *repositories.CardRepository,
	accountRepo *repositories.AccountRepository,
	pgpEntity *openpgp.Entity,
	logger *logrus.Logger,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
		pgpEntity:   pgpEntity,
		logger:      logger,
	}
}

func (s *CardService) GenerateCard(userID, accountID uint, cvv string) (*models.Card, error) {
	// Проверка прав доступа
	if _, err := s.accountRepo.GetByIDAndUser(accountID, userID); err != nil {
		return nil, fmt.Errorf("account access denied")
	}

	// Генерация данных карты
	cardNumber := generateLuhnValidNumber()
	expiry := time.Now().AddDate(5, 0, 0).Format("01/06")
	
	// Шифрование данных
	encryptedData, hmac, err := s.encryptCardData(fmt.Sprintf("%s|%s", cardNumber, expiry))
	if err != nil {
		return nil, err
	}

	// Хеширование CVV
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash CVV")
	}

	card := &models.Card{
		UserID:        userID,
		AccountID:     accountID,
		EncryptedData: encryptedData,
		Hmac:          hmac,
		CvvHash:       string(cvvHash),
	}

	if err := s.cardRepo.Create(card); err != nil {
		s.logger.Errorf("failed to save card: %v", err)
		return nil, fmt.Errorf("failed to save card: %w", err)
	}
	
	return card, nil
}

func (s *CardService) GetUserCards(userID uint) ([]models.Card, error) {
	return s.cardRepo.GetByUser(userID)
}

func generateLuhnValidNumber() string {
	rand.Seed(time.Now().UnixNano())
	number := "4" // Visa-подобные карты
	for i := 0; i < 14; i++ {
		number += string(byte(rand.Intn(10) + '0'))
	}
	
	sum := 0
	double := len(number)%2 == 0
	for _, c := range number {
		digit := int(c - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	
	checkDigit := (10 - (sum % 10)) % 10
	return number + fmt.Sprintf("%d", checkDigit)
}

func (s *CardService) encryptCardData(data string) (encrypted string, hmac string, err error) {
	// Шифрование PGP
	encrypted, err = crypto.EncryptPGP(data, s.pgpEntity)
	if err != nil {
		return "", "", fmt.Errorf("PGP encryption failed: %v", err)
	}

	// Генерация HMAC
	mac := hmacpkg.New(sha256.New, []byte(s.pgpEntity.PrivateKey.KeyIdString()))
	mac.Write([]byte(encrypted))
	hmac = hex.EncodeToString(mac.Sum(nil))

	return encrypted, hmac, nil
}