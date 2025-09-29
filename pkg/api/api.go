package api

import (
    "net/http"
    "time"
)

const DateFormat = "20060102"
const MaxDayInterval = 400

func Init() {
    http.HandleFunc("/api/nextdate", nextDayHandler)
    http.HandleFunc("/api/task", taskHandler)
    http.HandleFunc("/api/tasks", tasksHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    nowParam := r.URL.Query().Get("now")
    dateParam := r.URL.Query().Get("date")
    repeatParam := r.URL.Query().Get("repeat")

    var now time.Time
    if nowParam == "" {
        now = time.Now()
    } else {
        parsedNow, err := time.Parse(DateFormat, nowParam)
        if err != nil {
            http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
            return
        } 
        now = parsedNow
    }

    if dateParam == "" {
        http.Error(w, "Отсутствует параметр date", http.StatusBadRequest)
        return
    }
    if repeatParam == "" {
        http.Error(w, "Отсутствует параметр repeat", http.StatusBadRequest)
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