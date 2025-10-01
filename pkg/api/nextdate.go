package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// afterNow сравнивает только по дате, игнорируя время
func afterNow(date, now time.Time) bool {
	return date.Format(DateFormat) > now.Format(DateFormat)
}

// NextDate вычисляет следующую дату задачи по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", dstart)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid rule format")
	}

	ruleType := parts[0]

	// Если начальная дата больше now, возвращаем её
	if afterNow(startDate, now) {
		return startDate.Format(DateFormat), nil
	}

	switch ruleType {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d rule format: expected d <number>")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("interval must be a number")
		}
		if interval <= 0 || interval > MaxDayInterval {
			return "", errors.New("interval out of range")
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
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	default:
		return "", errors.New("unsupported rule format")
	}
}