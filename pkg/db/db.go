package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"sync"
	"time"
)

var (
	db   *sql.DB // Глобальная переменная подключения к БД
	dbMu sync.RWMutex // Мьютекс для защиты глобальной переменной
)

// schema содержит SQL-запросы для создания таблиц и индексов
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

// Init инициализирует подключение к базе данных и создает схему при необходимости
func Init(dbFile string) error {
	dbMu.Lock()
	defer dbMu.Unlock()

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

	// Закрываем существующее подключение если есть
	if db != nil {
		db.Close()
	}

	// Создаем новое подключение
	var openErr error
	db, openErr = sql.Open("sqlite", dbFile)
	if openErr != nil {
		return fmt.Errorf("failed to open database: %w", openErr)
	}

	// Настраиваем пул подключений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		db = nil
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Если файла не было, создаем схему
	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			db = nil
			return fmt.Errorf("failed to create schema: %w", err)
		}
		fmt.Printf("The database was created successfully: %s\n", dbFile)
	}

	return nil
}

// GetDB возвращает глобальное подключение к базе данных
func GetDB() *sql.DB {
	dbMu.RLock()
	defer dbMu.RUnlock()
	return db
}

// Close закрывает подключение к базе данных
func Close() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	
	if db != nil {
		err := db.Close()
		db = nil
		return err
	}
	return nil
}