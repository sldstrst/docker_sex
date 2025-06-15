package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Глобальная переменная для подключения к базе данных
var db *sql.DB

// initDB - инициализация подключения к базе данных и создание таблицы counters
func initDB() {
	var err error
	// Формируем строку подключения к базе данных из переменных окружения
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Пытаемся подключиться к базе с несколькими попытками и задержкой
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	// Если не удалось подключиться, завершаем программу с ошибкой
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Создаем таблицу counters, если она не существует, и инициализируем счетчик
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS counters (
			id SERIAL PRIMARY KEY,
			value INTEGER NOT NULL DEFAULT 0
		);
		
		INSERT INTO counters (value) 
		SELECT 0 
		WHERE NOT EXISTS (SELECT 1 FROM counters);
	`)
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
}

// main - точка входа в приложение
func main() {
	// Инициализируем базу данных
	initDB()
	// Закрываем подключение к базе при завершении работы
	defer db.Close()

	// Обработчик корневого пути, увеличивает счетчик и возвращает текущее значение
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Увеличиваем счетчик в базе
		_, err := db.Exec("UPDATE counters SET value = value + 1 WHERE id = 1 RETURNING value")
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Получаем текущее значение счетчика
		var count int
		err = db.QueryRow("SELECT value FROM counters WHERE id = 1").Scan(&count)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Возвращаем количество посещений в ответе
		fmt.Fprintf(w, "Total visits: %d", count)
	})

	// Логируем запуск сервера на порту 8081
	log.Println("Counter service with PostgreSQL started on :8081")
	// Запускаем HTTP сервер на порту 8081, логируем фатальные ошибки
	log.Fatal(http.ListenAndServe(":8081", nil))
}
