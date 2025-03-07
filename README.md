# Flexible Transfer Backend  

## 1. Описание и цель проекта  
Бэкенд-сервис для P2P платформы обмена валют между пользователями из разных стран. Основные цели:  
- Обеспечение кросс-валютных операций по принципу peer-to-peer  
- Автоматический подбор совпадающих ордеров с минимальной комиссией (0.5-1%)  
- Интеграция с внешними API для получения актуальных курсов валют  
- Механизм резервирования средств и подтверждения транзакций  
- Двухэтапная верификация пользователей (KYC)  

## 2. Архитектура системы  

### Основные компоненты  
```
├── delivery/      # HTTP обработчики и роутинг  
├── usecase/       # Бизнес-логика операций  
├── repository/    # Работа с БД и внешними сервисами  
├── domain/        # Ядро системы (сущности, валидация)  
└── config/        # Конфигурация приложения  
```

### Принцип работы
1. Пользователь создает заявку на обмен валюты через API
2. Система проверяет минимальный лимит и доступность валюты
3. Ищет совпадения среди активных ордеров с учетом отклонения ±5%
4. Резервирует средства на счетах участников
5. Выполняет обмен при подтверждении сделки через чат-интерфейс
6. Обновляет статусы ордеров и балансы пользователей
7. Сохраняет историю транзакций в реляционной БД

## 3. Требования
- Go 1.21+
- MySQL 8.0+
- Redis 6.2+ (для кеширования курсов валют)
- Docker 20.10+

## 4. Установка и запуск

### Шаг 1: Запуск инфраструктуры
```bash
docker-compose up -d mysql redis
```

### Шаг 2: Инициализация БД
```bash
mysql -h 127.0.0.1 -u root -p < init.sql
```

### Шаг 3: Настройка окружения  
Создайте .env файл:
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=flex_exchange
REDIS_URL=redis://localhost:6379
EXCHANGE_RATE_API_KEY=your_api_key
```

### Шаг 4: Сборка и запуск
```bash
go mod tidy
go build -o app ./cmd/server/
./app
```

## 5. Примеры API запросов

Создание ордера:
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order_123",
    "user_from": "user1",
    "currency_from": "USD",
    "currency_to": "EUR",
    "amount_from": 100,
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

Поиск совпадающих ордеров:
```bash
curl -X GET "http://localhost:8080/api/v1/orders/match?from=USD&to=EUR&amount=100"
```

## 6. Деплой

### Docker образ
```bash
docker build -t flexible-transfer-backend .
docker run -d --env-file .env -p 8080:8080 flexible-transfer-backend
```

### Kubernetes (пример конфигурации)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: transfer-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: transfer-backend
  template:
    metadata:
      labels:
        app: transfer-backend
    spec:
      containers:
      - name: backend
        image: flexible-transfer-backend:latest
        envFrom:
        - secretRef:
            name: backend-secrets
```

### Production настройки
- Используйте отдельного пользователя БД с ограниченными правами
- Настройте TLS для MySQL соединений
- Реализуйте механизм ротации логов
- Используйте Redis Sentinel для отказоустойчивости

## 7. Мониторинг
Журналирование операций осуществляется через Zap logger:  
```logs/application.log```  
Уровни логирования:
- DEBUG: Отладка операций
- INFO: Успешные транзакции  
- WARN: Подозрительные активности  
- ERROR: Критические ошибки

Метрики Prometheus доступны по endpoint `/metrics`