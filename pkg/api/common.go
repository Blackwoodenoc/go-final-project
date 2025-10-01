package api

import (
	"encoding/json"
	"net/http"
)

const (
	DateFormat     = "20060102" // Формат даты YYYYMMDD
	MaxDayInterval = 400        // Максимальный интервал в днях
)

// writeJSONSuccess отправляет успешный JSON ответ
func writeJSONSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}


// writeJSONError отправляет JSON ответ с ошибкой
func writeJSONError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": error})
}