# Базовый образ для сборки
FROM golang:1.23-alpine AS builder

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .
COPY .env .

# Собираем бинарник с оптимизацией
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/main.go

# Финальный минимальный образ
FROM alpine:latest

# Устанавливаем таймзону (опционально)
RUN apk --no-cache add tzdata curl

# Рабочая директория
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/bot .
COPY --from=builder /app/.env .

# Запускаем бота
CMD ["sh", "-c", "curl -f https://api.telegram.org || echo 'Network failed'; ./bot"]