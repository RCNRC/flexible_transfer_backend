# Flexible Transfer Backend

## 1. Описание и цель проекта
Бэкенд-сервис для P2P платформы обмена валют между пользователями из разных стран. Реализует децентрализованный обмен 
с минимальными комиссиями, автоматический подбор совпадающих ордеров и безопасное проведение транзакций. Основные цели:
- Обеспечение кросс-валютных операций между пользователями
- Поддержка актуальных курсов обмена через API
- Гарантия безопасности операций с двухфакторной проверкой
- Снижение комиссий до 0.5% за счет p2p модели

## 2. Архитектура и принцип работы

### Микросервисная архитектура
Проект реализован по принципам Clean Architecture:
- **Delivery layer**: HTTP handlers (`/internal/delivery`)
  - Роутинг запросов
  - Валидация входных данных
  - Работа с форматами JSON
- **UseCase layer**: Бизнес-логика (`/internal/core/usecase`)
  - Создание/обработка ордеров
  - Матчинг заявок
  - Валютные расчеты
- **Repository layer**: Работа с данными (`/internal/infrastructure/repository`)
  - MySQL: хранение ордеров/валют
  - Redis: кэширование курсов валют
- **Domain layer**: Ядро системы (`/internal/core/domain`)
  - Сущности бизнес-логики
  - Валидация правил

### Основные workflow
1. Пользователь создает ордер на обмен через API
2. Система валидирует баланс и минимальные лимиты
3. Поиск совпадающих ордеров по встречному курсу
4. Резервирование средств и подтверждение сделки
5. Обновление статусов ордеров в БД

## 3. Развертывание системы

### Требования
- Go 1.21+
- MySQL 8.0+ 
- Docker 20.10+
- Redis 6.2+

### Шаги деплоя

Команда 1 (Скачивание зависимостей)