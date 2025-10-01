package main

import (
    "fmt"
    "os"

    "go1f/pkg/db"
    "go1f/pkg/server"
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

    // Проверяем что БД инициализирована
    if db.GetDB() == nil {
        panic("Database is still nil after initialization!")
    }

    // Запускаем сервер
    if err := server.Run(); err != nil {
        panic(err)
    }
}