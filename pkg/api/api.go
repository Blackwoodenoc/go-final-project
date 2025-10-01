package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go1f/pkg/auth"
)

// SignInRequest структура для запроса входа
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse структура для ответа входа
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// Init регистрирует обработчики HTTP-запросов
func Init() {
	// Публичные маршруты (без аутентификации)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)

	// Защищенные маршруты (требуют аутентификации)
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(taskDoneHandler))
}

// authMiddleware проверяет аутентификацию
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если аутентификация отключена, пропускаем
		if !auth.IsAuthEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из куки
		var tokenString string
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}

		// Если токена нет в куках, проверяем заголовок Authorization
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		}

		// Проверяем токен
		valid, err := auth.ValidateToken(tokenString)
		if err != nil || !valid {
			writeJSONError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// signinHandler обрабатывает аутентификацию
func signinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	if r.Method != http.MethodPost {
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем пароль
	expectedPassword := getPassword()
	
	if req.Password != expectedPassword {
		writeJSONError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := auth.GenerateToken()
	if err != nil {
		writeJSONError(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	response := SignInResponse{Token: token}
	json.NewEncoder(w).Encode(response)
}

// Добавьте эту функцию
func getPassword() string {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "qwerty123" // тот же пароль по умолчанию
	}
	return password
}

// nextDayHandler обрабатывает запросы для вычисления следующей даты
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Invalid now parameter format", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	if dateParam == "" {
		http.Error(w, "Missing date parameter", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "Missing repeat parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
