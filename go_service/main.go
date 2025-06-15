package main

import (
	"fmt"
	"log"
	"net/http"
)

// main - точка входа в приложение
func main() {
	// Обработчик корневого пути, возвращает приветственное сообщение
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the simplest Go microservice!")
	})

	// Обработчик для проверки состояния сервиса (health check)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Возвращаем статус 200 OK
		fmt.Fprintf(w, "OK")         // Тело ответа "OK"
	})

	// Логируем запуск сервера на порту 8080
	log.Println("Starting server on :8080")
	// Запускаем HTTP сервер на порту 8080, логируем фатальные ошибки
	log.Fatal(http.ListenAndServe(":8080", nil))
}
