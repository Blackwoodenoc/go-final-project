package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Task представляет задачу в планировщике
type Task struct {
    ID      string `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных
// Возвращает ID добавленной задачи или ошибку
func AddTask(task *Task) (int64, error) {
    if GetDB() == nil {
        return 0, fmt.Errorf("database connection is not initialized")
    }
    
    var id int64
    query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES(?, ?, ?, ?)`
    res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
    if err == nil {
        id, err = res.LastInsertId()
    }
    return id, err
}

// Tasks возвращает список задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {
    if GetDB() == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    rows, err := GetDB().Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
    if err != nil {
        return nil, fmt.Errorf("database query error: %w", err)
    }
    defer rows.Close()

    return scanTasks(rows)
}

// SearchTasks ищет задачи по подстроке или дате
func SearchTasks(search string, limit int) ([]*Task, error) {
    if GetDB() == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    // Проверяем, является ли search датой в формате 02.01.2006
    if date, ok := parseDate(search); ok {
        // Поиск по дате
        rows, err := GetDB().Query(
            "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit",
            sql.Named("date", date),
            sql.Named("limit", limit),
        )
        if err != nil {
            return nil, fmt.Errorf("database query error: %w", err)
        }
        defer rows.Close()
        
        return scanTasks(rows)
    }
    
    // Поиск по подстроке в заголовке или комментарии
    searchPattern := "%" + strings.ToLower(search) + "%"
    rows, err := GetDB().Query(
        `SELECT id, date, title, comment, repeat FROM scheduler 
         WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search 
         ORDER BY date LIMIT :limit`,
        sql.Named("search", searchPattern),
        sql.Named("search", searchPattern),
        sql.Named("limit", limit),
    )
    if err != nil {
        return nil, fmt.Errorf("database query error: %w", err)
    }
    defer rows.Close()
    
    return scanTasks(rows)
}

// parseDate пытается разобрать строку как дату в формате 02.01.2006
// Возвращает дату в формате 20060102 и true, если разбор успешен
func parseDate(dateStr string) (string, bool) {
    // Убираем лишние пробелы
    dateStr = strings.TrimSpace(dateStr)
    
    // Пытаемся разобрать дату в формате 02.01.2006
    parsedDate, err := time.Parse("02.01.2006", dateStr)
    if err != nil {
        return "", false
    }
    
    // Преобразуем в формат 20060102
    return parsedDate.Format("20060102"), true
}

// scanTasks сканирует строки из результата запроса и возвращает список задач
func scanTasks(rows *sql.Rows) ([]*Task, error) {
    var tasks []*Task
    
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }
        tasks = append(tasks, &task)
    }
    
    // Проверяем ошибки после итерации
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    
    // Если tasks равен nil, возвращаем пустой слайс
    if tasks == nil {
        tasks = []*Task{}
    }
    
    return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
    if GetDB() == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    var task Task
    err := GetDB().QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
        Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("task not found")
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
    if GetDB() == nil {
        return fmt.Errorf("database connection is not initialized")
    }

    query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
    res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return fmt.Errorf("update error: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("update check error: %w", err)
    }
    
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
    if GetDB() == nil {
        return fmt.Errorf("database connection is not initialized")
    }

    res, err := GetDB().Exec("DELETE FROM scheduler WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("delete error: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("delete check error: %w", err)
    }
    
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

// UpdateDate обновляет дату задачи
func UpdateDate(next string, id string) error {
    if GetDB() == nil {
        return fmt.Errorf("database connection is not initialized")
    }

    query := `UPDATE scheduler SET date = ? WHERE id = ?`
    res, err := GetDB().Exec(query, next, id)
    if err != nil {
        return fmt.Errorf("update error: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("update check error: %w", err)
    }
    
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}