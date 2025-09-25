package main

import (
	"fmt"
	"go1f/pkg/server"
	"go1f/pkg/db"
	"os"
)

func main() {
	// Определяем путь к БД
	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	// Инициализируем БД
	if err := db.Init(dbFile); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer db.Close()

	// Запускаем сервер
	if err := server.Run(); err != nil {
		panic(err)
	}
}