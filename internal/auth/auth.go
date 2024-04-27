package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int
	jwt.StandardClaims
}

var jwtKey = []byte("1234")

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsPasswordSafe(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func GetUserIdFromToken(tokenString string) (int, error) {
    token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
    if err != nil {
        fmt.Println("Error during ParseUnverified:", err)
        return 0, fmt.Errorf("failed to parse token: %s", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return 0, errors.New("invalid claims")
    }

    userID, err := strconv.Atoi(fmt.Sprintf("%v", claims.UserID))
    if err != nil {
        return 0, fmt.Errorf("failed to parse user ID: %s", err)
    }

    return userID, nil
}

func AuthenticateUser(r *http.Request) (int, error) {
    bearerToken := r.Header.Get("Authorization")
    if bearerToken == "" {
        return 0, errors.New("authorization token required")
    }

    token := strings.Split(bearerToken, " ")
    if len(token) != 2 {
        return 0, errors.New("invalid authorization token")
    }

    userID, err := GetUserIdFromToken(token[1])
    if err != nil {
        return 0, errors.New("failed to authenticate user")
    }

    return userID, nil
}
