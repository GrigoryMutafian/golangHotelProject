# ---------- Этап 1: сборка ----------
FROM golang:1.25.2-alpine AS builder

# Рабочая директория внутри контейнера
WORKDIR /app

# Скопируем go.mod и go.sum и заранее выкачаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопируем весь исходный код
COPY . .

# Собираем бинарь (имя: server)
RUN go build -o server .
# ---------- Этап 2: запуск ----------
    
FROM alpine:3.20

WORKDIR /app

# Копируем бинарь из builder
COPY --from=builder /app/server .

# (необязательно, но удобно: добавить сертификаты для https-запросов)
RUN apk add --no-cache ca-certificates

# Открываем порт (например 8080)
EXPOSE 8080

# Запуск бинаря
CMD ["./server"]
