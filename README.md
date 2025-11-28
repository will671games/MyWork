# Wallet Service

REST API сервис для управления кошельками и транзакциями.

## Описание

Приложение предоставляет API для создания кошельков, выполнения транзакций (депозиты и выводы) и получения баланса. Сервис разработан с учетом работы в конкурентной среде и обеспечивает целостность данных при высоких нагрузках.

## Технологический стек

- **Golang** - язык программирования
- **PostgreSQL** - база данных
- **Docker** - контейнеризация
- **Echo** - веб-фреймворк
- **GORM** - ORM для работы с БД
- **go.uber.org/fx** - dependency injection

## API Endpoints

### Создание кошелька

```
POST /api/v1/wallets
Content-Type: application/json

{
  "balance": 1000.0
}
```

**Response:**
```json
{
  "walletId": "123e4567-e89b-12d3-a456-426614174000",
  "balance": 1000.0
}
```

### Получение баланса

```
GET /api/v1/wallets/{walletId}
```

**Response:**
```json
{
  "walletId": "123e4567-e89b-12d3-a456-426614174000",
  "balance": 1500.0
}
```

### Выполнение транзакции

```
POST /api/v1/wallet
Content-Type: application/json

{
  "wallet_id": "123e4567-e89b-12d3-a456-426614174000",
  "operation_type": "DEPOSIT",
  "amount": 500.0
}
```

**Параметры:**
- `wallet_id` (string, UUID, required) - идентификатор кошелька
- `operation_type` (string, required) - тип операции: `DEPOSIT` или `WITHDRAW`
- `amount` (float, required, > 0) - сумма транзакции

**Response:**
```json
{
  "wallet_id": "123e4567-e89b-12d3-a456-426614174000",
  "operation_type": "DEPOSIT",
  "amount": 500.0,
  "balance": 1500.0
}
```

**Ошибки:**
- `400 Bad Request` - невалидные данные или недостаточно средств
- `404 Not Found` - кошелек не найден
- `500 Internal Server Error` - внутренняя ошибка сервера

## Структура проекта

```
TestProject/
├── source/
│   ├── cmd/
│   │   └── main.go              # Точка входа приложения
│   ├── config/
│   │   └── config.go             # Конфигурация
│   └── internal/
│       ├── application/         # Бизнес-логика
│       ├── entities/             # Доменные сущности
│       ├── storage/              # Слой работы с БД
│       └── transport/            # HTTP handlers
├── docker-compose.yml
├── config.env.example
└── README.md
```

## Примеры использования

### Создание кошелька

```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"balance": 1000.0}'
```

### Получение баланса

```bash
curl http://localhost:8080/api/v1/wallets/{walletId}
```

### Депозит

```bash
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": "123e4567-e89b-12d3-a456-426614174000",
    "operation_type": "DEPOSIT",
    "amount": 500.0
  }'
```

### Вывод средств

```bash
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": "123e4567-e89b-12d3-a456-426614174000",
    "operation_type": "WITHDRAW",
    "amount": 300.0
  }'
```
