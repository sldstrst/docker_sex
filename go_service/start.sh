#!/bin/sh

# Запускаем оба сервиса в фоне
./greeting-service &
pid_greeting=$!  # PID процесса greeting-service

./counter-service &
pid_counter=$!  # PID процесса counter-service

# Функция для обработки сигналов завершения (TERM, INT)
term_handler() {
  echo "Stopping services..."
  # Отправляем сигнал завершения процессам сервисов
  kill -TERM "$pid_greeting" 2>/dev/null
  kill -TERM "$pid_counter" 2>/dev/null
  # Ожидаем завершения процессов
  wait "$pid_greeting"
  wait "$pid_counter"
  exit 0
}

# Устанавливаем обработчик сигналов TERM и INT
trap term_handler TERM INT

# Ждем завершения процессов сервисов, чтобы контейнер не завершился
wait "$pid_greeting"
wait "$pid_counter"
