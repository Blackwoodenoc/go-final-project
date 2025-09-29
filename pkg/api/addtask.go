package api

import (
    "encoding/json"
    "net/http"
    "time"
    "go1f/pkg/db"
)


func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        addTaskHandler(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    
    var task db.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        writeJSONError(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Проверяем обязательное поле title
    if task.Title == "" {
        writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
        return
    }

    // Обрабатываем дату
    if err := processTaskDate(&task); err != nil {
        writeJSONError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Добавляем задачу в базу
    id, err := db.AddTask(&task)
    if err != nil {
        writeJSONError(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSONSuccess(w, map[string]interface{}{"id": id})
}

func processTaskDate(task *db.Task) error {
    now := time.Now()
    
    // Если дата не указана, используем сегодняшнюю
    if task.Date == "" {
        task.Date = now.Format(DateFormat)
    }

    // Проверяем формат даты
    t, err := time.Parse(DateFormat, task.Date)
    if err != nil {
        return err
    }

    // Если дата в прошлом
    if t.Before(now) {
        if task.Repeat == "" {
            // Без повторения - используем сегодня
            task.Date = now.Format(DateFormat)
        } else {
            // С повторением - вычисляем следующую дату
            next, err := NextDate(now, task.Date, task.Repeat)
            if err != nil {
                return err
            }
            task.Date = next
        }
    }

    return nil
}

func writeJSONSuccess(w http.ResponseWriter, data map[string]interface{}) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}

func writeJson(w http.ResponseWriter, data interface{}) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, error string, code int) {
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": error})
}

