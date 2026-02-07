API для управления данными отеля
Небольшой учебный проект на Go для управления номерами и бронированиями в отеле. Сделан с упором на Clean Architecture. Бизнес-логика отделена от хранилища и HTTP-слоя.
Что умеет:
  -Номера: создать, обновить, удалить, получить список (в том числе с фильтрами).
  -Бронирования: создать, обновить, удалить, получить по ID, получить список (в том числе с фильтрами).

Быстро развернуть проект с помощью Docker Compose командой: (bash) "docker-compose up --build"

Сервис слушает порт :8080

Swagger-документация:

  После запуска приложения, открыть в браузере http://localhost:8080/swagger/index.html

Rooms
  POST /Create — создать номер

  PATCH /Patch?id=... — обновить номер

  DELETE /RemoveRoom — удалить номер

  GET /GetFilteredRooms — получить список номеров

Bookings
  POST /CreateBooking — создать бронирование

  GET /ReadBookingByID?id=... — получить бронирование по ID

  PATCH /PatchBookingByID — обновить бронирование

  GET /GetFilteredBookings — список бронирований

  DELETE /RemoveBooking — удалить бронирование


Переменные окружения для БД:

DB_HOST,DB_PORT,DB_USER,DB_PASSWORD,DB_NAME