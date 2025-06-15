# Используем официальный образ Golang на базе Alpine для сборки counter-service
FROM golang:1.24-alpine as builder

WORKDIR /app

# Копируем go.mod для кэширования зависимостей
COPY go.mod ./
RUN go mod tidy

# Копируем исходники counter-service
COPY ./counter ./counter

WORKDIR /app/counter

# Собираем counter-service
RUN go build -o counter-service .

# Финальный образ
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache libc6-compat

# Добавляем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/counter/counter-service .

EXPOSE 8081

USER appuser

CMD ["./counter-service"]
