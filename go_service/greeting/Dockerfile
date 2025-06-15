# Используем официальный образ Golang на базе Alpine для сборки greeting-service
FROM golang:1.24-alpine as builder

WORKDIR /app

# Копируем go.mod для кэширования зависимостей
COPY go.mod ./
RUN go mod tidy

# Копируем исходники greeting-service
COPY . .

# Собираем greeting-service
RUN go build -o greeting-service .

# Финальный образ
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache libc6-compat

# Добавляем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/greeting-service .

EXPOSE 8080

USER appuser

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget -q --spider http://localhost:8080/health || exit 1

CMD ["./greeting-service"]
