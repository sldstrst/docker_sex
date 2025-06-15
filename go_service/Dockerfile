# Используем официальный образ Golang на базе Alpine для сборки приложения
FROM golang:1.24-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файл go.mod для кэширования зависимостей
COPY go.mod ./
# Загружаем зависимости, чтобы кэшировать их слой
RUN go mod tidy

# Копируем все файлы проекта в контейнер
COPY . .

# Собираем основной сервис greeting-service
RUN go build -o greeting-service .
# Переходим в директорию второго сервиса counter
WORKDIR /app/counter
# Собираем сервис counter-service
RUN go build -o counter-service .

# Используем минимальный образ Alpine для финального контейнера
FROM alpine:latest
# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем необходимые библиотеки для запуска приложений
RUN apk add --no-cache libc6-compat bash

# Добавляем нового пользователя и группу для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем собранные бинарники из builder-образа
COPY --from=builder /app/greeting-service .
COPY --from=builder /app/counter/counter-service .
# Копируем скрипт запуска
COPY --from=builder /app/start.sh .

# Делаем скрипт запуска исполняемым
RUN chmod +x ./start.sh

# Открываем порты, на которых работают сервисы
EXPOSE 8080
EXPOSE 8081

# Запускаем контейнер от имени созданного пользователя
USER appuser

# Добавляем проверку здоровья контейнера, проверяя доступность /health на порту 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget -q --spider http://localhost:8080/health || exit 1

# Запускаем скрипт старта, который запускает оба сервиса
CMD ["./start.sh"]
