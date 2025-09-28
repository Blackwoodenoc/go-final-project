package api

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"
)

// Сравнение только по дате, игнорируем время
func afterNow(date, now time.Time) bool {
    return date.After(now)
}

// NextDate вычисляет следующую дату задачи по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
    if repeat == "" {
        return "", errors.New("пустое правило повторения")
    }

    startDate, err := time.Parse(DateFormat, dstart)
    if err != nil {
        return "", fmt.Errorf("некорректный формат даты: %s", dstart)
    }

    parts := strings.Split(repeat, " ")
    if len(parts) == 0 {
        return "", errors.New("неверный формат правила")
    }

    ruleType := parts[0]

    // Если начальная дата больше now, возвращаем её
    if afterNow(startDate, now) {
        return startDate.Format(DateFormat), nil
    }

    switch ruleType {
    case "d":
        if len(parts) != 2 {
            return "", errors.New("неверный формат правила d: ожидается d <число>")
        }
        interval, err := strconv.Atoi(parts[1])
        if err != nil {
            return "", errors.New("интервал должен быть числом")
        }
        if interval <= 0 || interval > MaxDayInterval {
            return "", nil // возвращаем пустую строку, если интервал больше допустимого
        }

        date := startDate
        for {
            date = date.AddDate(0, 0, interval)
            if afterNow(date, now) {
                return date.Format(DateFormat), nil
            }
        }

    case "y":
    date := startDate
    for {
    // Если дата 29 февраля, следующий год не високосный → ставим 1 марта
        if date.Day() == 29 && date.Month() == 2 {
            nextYear := date.Year() + 1
        if !isLeap(nextYear) {
            date = time.Date(nextYear, 3, 1, 0, 0, 0, 0, date.Location())
        } else {
            date = time.Date(nextYear, 2, 29, 0, 0, 0, 0, date.Location())
        }
    } else {
        date = date.AddDate(1, 0, 0)
    }
    if afterNow(date, now) {
        return date.Format(DateFormat), nil
    }
    }

    default:
        return "", errors.New("неподдерживаемый формат правила")
    }
}

// Проверка високосного года
func isLeap(year int) bool {
    return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}