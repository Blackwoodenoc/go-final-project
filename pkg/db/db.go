package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

func Init(dbFile string) error {
	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("failed to check database file: %w", err)
		}
	}

	// Открываем базу данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Если файла не было, создаем схему
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
		fmt.Printf("Database created successfully: %s\n", dbFile)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}