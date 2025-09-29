// package db

// type Task struct {
//     ID      string `json:"id"`
//     Date    string `json:"date"`
// 	Title 	string `json:"title"`
// 	Comment string `json:"comment"`
// 	Repeat  string `json:"repeat"`
// }

// func AddTask(task *Task) (int64, error) {
//     var id int64
//     // определяем запрос
//     query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES(?, ?, ?, ?)`
//     res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
//     if err == nil {
//         id, err = res.LastInsertId()
//     }
//     return id, err
// }

package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
    ID      string `json:"id"`
    Date    string `json:"date"`
	Title 	string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
    // Добавляем проверку на nil
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

func Tasks(limit int) ([]*Task, error) {
    if GetDB() == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }

    rows, err := GetDB().Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
    if err != nil {
        return nil, fmt.Errorf("database query error: %w", err)
    }
    defer rows.Close()

    var tasks []*Task

    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }
        tasks = append(tasks, &task)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    // Возвращаем пустой слайс вместо nil
    if tasks == nil {
        tasks = []*Task{}
    }

    return tasks, nil
}