package api

import (
	"encoding/json"
	"net/http"
	"time"
	"go1f/pkg/db"
)

// taskHandler обрабатывает все методы для работы с задачей
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete: 
		deleteTaskHandler(w, r)
	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getTaskHandler обрабатывает получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "ID not specified", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSONSuccess(w, task, http.StatusOK)
}

// updateTaskHandler обрабатывает обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSONError(w, "Task title not specified", http.StatusBadRequest)
		return
	}

	// Проверяем наличие ID
	if task.ID == "" {
		writeJSONError(w, "Task ID not specified", http.StatusBadRequest)
		return
	}

	// Обрабатываем дату
	if err := processTaskDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе
	if err := db.UpdateTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSONSuccess(w, map[string]interface{}{}, http.StatusOK)
}

// addTaskHandler обрабатывает создание новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSONError(w, "Task title not specified", http.StatusBadRequest)
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
		writeJSONError(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, map[string]interface{}{"id": id}, http.StatusOK)
}

// processTaskDate обрабатывает и валидирует дату задачи
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

// taskDoneHandler обрабатывает отметку о выполнении задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "ID parameter missing", http.StatusBadRequest)
		return
	}
	
	// Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Если задача без повтора - удаляем
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONSuccess(w, map[string]interface{}{}, http.StatusOK)
		return
	}
	
	// Для периодической задачи вычисляем следующую дату
	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Обновляем дату
	err = db.UpdateDate(nextDate, id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, map[string]interface{}{}, http.StatusOK)
}

// deleteTaskHandler обрабатывает удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "ID not specified", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, map[string]interface{}{}, http.StatusOK)
}