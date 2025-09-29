// package db

// import (
// 	"database/sql"
// 	"fmt"
// 	_ "modernc.org/sqlite"
// 	"os"
// )

// var db *sql.DB

// const schema = `
// CREATE TABLE scheduler (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     date CHAR(8) NOT NULL DEFAULT "",
//     title VARCHAR(128) NOT NULL,
//     comment TEXT NOT NULL,
//     repeat VARCHAR(128) NOT NULL
// );

// CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

// func Init(dbFile string) error {
// 	// Проверяем существование файла БД
// 	_, err := os.Stat(dbFile)

// 	var install bool
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			install = true
// 		} else {
// 			return fmt.Errorf("не удалось проверить файл базы данных: %w", err)
// 		}
// 	}

// 	// Открываем базу данных
// 	db, err := sql.Open("sqlite", dbFile)
// 	if err != nil {
// 		return fmt.Errorf("не удалось открыть базу данных: %w", err)
// 	}

// 	// Проверяем соединение
// 	if err := db.Ping(); err != nil {
// 		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
// 	}

// 	// Если файла не было, создаем схему
// 	if install {
// 		if _, err := db.Exec(schema); err != nil {
// 			return fmt.Errorf("не удалось создать схему: %w", err)
// 		}
// 		fmt.Printf("База данных успешно создана: %s\n", dbFile)
// 	}

// 	return nil
// }

// func GetDB() *sql.DB {
// 	return db
// }

// func Close() error {
// 	if db != nil {
// 		return db.Close()
// 	}
// 	return nil
// }


package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB  // Глобальная переменная

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
			return fmt.Errorf("не удалось проверить файл базы данных: %w", err)
		}
	}

	// ОТКЛЮЧАЕМ ЛОКАЛЬНУЮ ПЕРЕМЕННУЮ - используем глобальную db
	var openErr error
	db, openErr = sql.Open("sqlite", dbFile)  // Используем глобальную db
	if openErr != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", openErr)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Если файла не было, создаем схему
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("не удалось создать схему: %w", err)
		}
		fmt.Printf("База данных успешно создана: %s\n", dbFile)
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