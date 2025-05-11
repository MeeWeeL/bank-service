# Bank Service

Этот проект представляет собой банковский сервис, реализованный на языке Go с использованием PostgreSQL для хранения данных. Сервис поддерживает операции с пользователями, счетами, транзакциями и кредитами.

## Структура проекта

├── .env
├── Dockerfile
├── go.mod
├── migrations
│ ├── 001_init.up.sql
│ ├── 001_init.down.sql
│ ├── 002_accounts.up.sql
│ ├── 002_accounts.down.sql
│ ├── 003_cards.up.sql
│ ├── 003_cards.down.sql
│ ├── 004_transactions.up.sql
│ ├── 004_transactions.down.sql
│ ├── 005_payment_schedules.up.sql
│ ├── 005_payment_schedules.down.sql
│ ├── 006_credits.up.sql
│ └── 006_credits.down.sql
└── src
└── main.go


## Установка

1. **Клонируйте репозиторий:**

   ```bash
   git clone <URL-репозитория>
   cd <имя-репозитория>

    Настройте переменные окружения:

    Создайте файл .env в корневой директории проекта и добавьте следующие строки:

    DATABASE_URL=postgres://bank_user:bank_password@localhost:5432/bank_db?sslmode=disable
    JWT_SECRET=very_secret_key

    Соберите проект с помощью Docker:

    docker build -t bank-service .

    Запустите контейнер:

    docker run -p 8080:8080 --env-file .env bank-service

## Как пользоваться сервисом

Сервис предоставляет API для выполнения различных операций. Вот основные команды и их описание:
Пользователи

    Регистрация пользователя:
        POST /api/users/register
    Авторизация пользователя:
        POST /api/users/login
    Получение информации о пользователе:
        GET /api/users/{id}
    Получение списка всех пользователей:
        GET /api/users

Счета

    Создание счета:
        POST /api/accounts
    Получение информации о счете:
        GET /api/accounts/{id}
    Получение всех счетов пользователя:
        GET /api/accounts/user/{userId}
    Обновление счета:
        PUT /api/accounts/{id}
    Удаление счета:
        DELETE /api/accounts/{id}

Транзакции

    Создание транзакции:
        POST /api/transactions
    Получение информации о транзакции:
        GET /api/transactions/{id}
    Получение всех транзакций счета:
        GET /api/transactions/account/{accountId}

Кредиты

    Создание кредита:
        POST /api/credits
    Получение информации о кредите:
        GET /api/credits/{id}
    Получение всех кредитов пользователя:
        GET /api/credits/user/{userId}
    Обновление кредита:
        PUT /api/credits/{id}
    Удаление кредита:
        DELETE /api/credits/{id}

Как протестировать код

    Запустите тесты:

    В директории src выполните команду:

    go test ./...

    Проверьте миграции базы данных:

    Убедитесь, что все миграции применены:

    psql -U bank_user -d bank_db -f migrations/001_init.up.sql

    Замените имя файла на нужное для применения других миграций.

Заключение

Этот проект является демонстрацией банковского сервиса и может быть расширен дополнительными функциями. Если у вас есть вопросы или предложения, не стесняйтесь обращаться.