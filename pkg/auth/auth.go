package auth

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secretKey будет вычисляться на основе пароля
func getSecretKey() []byte {
	password := getPassword()
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func getPassword() string {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "qwerty123" // пароль по умолчанию
	}
	return password
}

// Claims структура для JWT claims
type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT токен
func GenerateToken() (string, error) {
	if !IsAuthEnabled() {
		return "no_auth", nil
	}

	// Создаем claims с информацией
	claims := &Claims{
		Login: "todo_user", 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "todo_app",
		},
	}

	// Создаем токен с claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Подписываем токен ключом на основе пароля
	signedToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) (bool, error) {
	if !IsAuthEnabled() {
		return true, nil
	}

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// IsAuthEnabled проверяет включена ли аутентификация
func IsAuthEnabled() bool {
	return getPassword() != ""
}