# Golang Hotel Project

Проект для управления отелем, написанный на Go с использованием архитектуры Clean Architecture.

## Описание проекта

Этот проект предоставляет RESTful API для управления отелем, включая управление номерами и бронированиями. Проект следует принципам Clean Architecture, разделяя бизнес-логику от доступа к данным.

## Структура проекта

```
golangHotelProject/
├── main.go                 # Точка входа в приложение
├── go.mod                  # Файл зависимостей Go
├── docker-compose.yml      # Конфигурация Docker Compose
├── Dockerfile              # Конфигурация Docker
├── init.sql               # SQL скрипт для инициализации базы данных
├── docs/                  # Документация Swagger
│   ├── docs.go            # Сгенерированная документация
│   ├── swagger.json       # JSON файл спецификации Swagger
│   └── swagger.yaml       # YAML файл спецификации Swagger
├── internal/              # Внутренняя структура проекта
│   ├── model/             # Модели данных
│   │   ├── bookingModel.go
│   │   └── hotelModel.go
│   ├── repository/        # Слой доступа к данным
│   │   ├── room.go
│   │   ├── booking.go
│   │   └── db/
│   │       └── db.go
│   ├── usecase/           # Бизнес-логика
│   │   ├── room.go
│   │   ├── booking.go
│   │   └── booking_test.go
│   ├── delivery/          # Слой доставки (handlers)
│   │   ├── handlers/
│   │   │   ├── room.go
│   │   │   ├── booking.go
│   │   │   └── dto/
│   │   │       └── dto.go
│   │   └── middleware/
│   │       └── cors.go
```

## Требования

- Go 1.19 или выше
- Docker и Docker Compose (для запуска в контейнерах)
- PostgreSQL (база данных)

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/GrigoryMutafian/golangHotelProject.git
cd golangHotelProject
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Настройка базы данных

1. Убедитесь, что PostgreSQL запущен
2. Выполните SQL скрипт для инициализации базы данных:

```bash
psql -h localhost -U your_username -d your_database -f init.sql
```

### 4. Запуск приложения

#### Запуск напрямую

```bash
go run main.go
```

#### Запуск через Docker Compose

```bash
docker-compose up --build
```

Приложение будет доступно на порту 8080.

## API Документация

После запуска приложения вы можете получить доступ к Swagger UI по адресу:

```
http://localhost:8080/swagger/index.html
```

### API Эндпоинты

#### Управление номерами

- **POST /Create** - Создание нового номера
  - Тело запроса: `{"roomType": "standard", "price": 150.00, "roomCount": 1, "floor": 1, "sleepingPlaces": 2}`
  
- **PATCH /Patch** - Обновление существующего номера
  - Параметры: `id` (query)
  - Тело запроса: `{"roomCount": 2, "floor": 2}`
  
- **DELETE /RemoveRoom** - Удаление номера
  - Тело запроса: `{"roomId": 1}`
  
- **POST /GetFilteredRooms** - Получение списка номеров с фильтрацией
  - Тело запроса: `{"roomType": "luxury", "floor": 2}` (опционально)

#### Управление бронированиями

- **POST /CreateBooking** - Создание нового бронирования
  - Тело запроса: `{"roomId": 1, "guestId": 1, "startDate": "2024-01-15T00:00:00Z", "endDate": "2024-01-20T00:00:00Z", "status": "confirmed"}`
  
- **GET /ReadBookingByID** - Получение бронирования по ID
  - Параметры: `id` (query)
  
- **PATCH /PatchBookingByID** - Обновление существующего бронирования
  - Тело запроса: `{"id": 1, "status": "cancelled"}`
  
- **GET /GetFilteredBookings** - Получение списка бронирований с фильтрацией
  - Тело запроса: `{"roomId": 1}` (опционально)
  
- **DELETE /RemoveBooking** - Удаление бронирования
  - Тело запроса: `1`

## Модели данных

### Room (Номер)
- `ID` - Уникальный идентификатор
- `RoomType` - Тип номера (standard, luxury, etc.)
- `Price` - Цена за ночь
- `RoomCount` - Количество номеров
- `Floor` - Этаж
- `SleepingPlaces` - Количество спальных мест
- `IsOccupied` - Занятость
- `NeedCleaning` - Требуется уборка

### Booking (Бронирование)
- `ID` - Уникальный идентификатор
- `RoomID` - ID номера
- `GuestID` - ID гостя
- `StartDate` - Дата начала бронирования
- `EndDate` - Дата окончания бронирования
- `Status` - Статус бронирования (confirmed, cancelled, etc.)

## Тестирование

Для запуска тестов:

```bash
go test ./...
```

## Конфигурация

Приложение использует следующие переменные окружения:

- `DB_HOST` - Хост базы данных
- `DB_PORT` - Порт базы данных
- `DB_USER` - Пользователь базы данных
- `DB_PASSWORD` - Пароль базы данных
- `DB_NAME` - Имя базы данных

## Лицензия

Этот проект распространяется под лицензией MIT.

## Контакты

Для вопросов и предложений, пожалуйста, обращайтесь к автору проекта.