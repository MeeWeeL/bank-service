package main

import (
	"bank-service/src/config"
	"bank-service/src/handlers"
	"bank-service/src/middleware"
	"bank-service/src/repositories"
	"bank-service/src/services"
	"bank-service/src/crypto"
	"database/sql"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	cfg := config.Load()

	// Инициализация БД
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db, logger)
	accountRepo := repositories.NewAccountRepository(db, logger)
	cardRepo := repositories.NewCardRepository(db, logger)
	transactionRepo := repositories.NewTransactionRepository(db, logger)
	creditRepo := repositories.NewCreditRepository(db, logger)
	paymentScheduleRepo := repositories.NewPaymentScheduleRepository(db, logger)
	

	// Инициализация PGP
	pgpEntity, err := crypto.InitPGP()
	if err != nil {
		logger.Fatal("Failed to initialize PGP: ", err)
	}

	// Инициализация сервисов
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, logger)
	accountService := services.NewAccountService(accountRepo, transactionRepo, logger)
	cardService := services.NewCardService(
		cardRepo, 
		accountRepo, 
		pgpEntity,
		logger,
	)
	cbrService := services.NewCBRService()
	emailService := services.NewEmailService(
		cfg.EmailHost,
		cfg.EmailPort,
		cfg.EmailUser,
		cfg.EmailPass,
		cfg.EmailFrom,
		logger,
	)
	creditService := services.NewCreditService(
		creditRepo, 
		paymentScheduleRepo, 
		accountService, 
		cbrService,
		logger,
		userRepo,
		emailService,
	)
	analyticsService := services.NewAnalyticsService(
		transactionRepo, 
		creditRepo, 
		accountRepo,
		paymentScheduleRepo,
	)

    go func() {
        ticker := time.NewTicker(12 * time.Hour)
        for range ticker.C {
            if err := creditService.ProcessOverduePayments(); err != nil {
                logger.Errorf("Overdue payments processing failed: %v", err)
            }
        }
    }()

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService, logger)
	accountHandler := handlers.NewAccountHandler(accountService, analyticsService, logger)
	transferHandler := handlers.NewTransferHandler(accountService, logger)
	creditHandler := handlers.NewCreditHandler(creditService, logger)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, logger)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, logger)

	// Инициализация сервиса карт
	cardHandler := handlers.NewCardHandler(cardService, logger)

	// Настройка маршрутизатора
	router := mux.NewRouter()

	// Публичные маршруты
	public := router.PathPrefix("/").Subrouter()
	public.HandleFunc("/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Защищенные маршруты
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.Handle)
	
	// Маршруты для счетов
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts/{accountId}", accountHandler.GetAccount).Methods("GET")
	protected.HandleFunc("/accounts/{accountId}/predict", accountHandler.PredictBalance).Methods("GET")

	// Для карт
	protected.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	protected.HandleFunc("/cards", cardHandler.GetCards).Methods("GET")

	// Трансферы
	protected.HandleFunc("/transfer", transferHandler.Transfer).Methods("POST")
	protected.HandleFunc("/accounts/{accountId}/deposit", accountHandler.Deposit).Methods("PUT")

	// Кредиты
	protected.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	protected.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetSchedule).Methods("GET")
	protected.HandleFunc("/credits/{accountId}/credits", creditHandler.GetCreditsByAccount).Methods("GET")

	// Аналитика
	protected.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")

	// Тестовый маршрут
	protected.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected route"))
	}).Methods("GET")

	// Запуск сервера
	logger.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal("Server failed: ", err)
	}
}