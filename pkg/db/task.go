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

import "fmt"
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