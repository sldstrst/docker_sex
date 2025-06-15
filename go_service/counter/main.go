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

var db *sql.DB

func initDB() {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Пытаемся подключиться с таймаутами
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

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Создаем таблицу если не существует
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

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Увеличиваем счетчик
		_, err := db.Exec("UPDATE counters SET value = value + 1 WHERE id = 1 RETURNING value")
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Получаем текущее значение
		var count int
		err = db.QueryRow("SELECT value FROM counters WHERE id = 1").Scan(&count)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Total visits: %d", count)
	})

	log.Println("Counter service with PostgreSQL started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}