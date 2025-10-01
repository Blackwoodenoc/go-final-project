package api

import (
	"net/http"
	"go1f/pkg/db"
)

// TasksResp структура для ответа с задачами
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает запросы на получение задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из query string
	search := r.URL.Query().Get("search")
	
	// Парсим лимит, по умолчанию 50
	limit := 50

	var tasks []*db.Task
	var err error
	
	if search != "" {
		// Если есть параметр search, используем поиск
		tasks, err = db.SearchTasks(search, limit)
	} else {
		// Иначе получаем все задачи
		tasks, err = db.Tasks(limit)
	}
	
	if err != nil {
		writeJSONError(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}