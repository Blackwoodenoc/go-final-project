package api

import (
    "net/http"
    "go1f/pkg/db"
)

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    
    // Добавьте проверку метода
    if r.Method != http.MethodGet {
        writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    tasks, err := db.Tasks(50)
    if err != nil {
        writeJSONError(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    writeJson(w, TasksResp{
        Tasks: tasks,
    })
}
