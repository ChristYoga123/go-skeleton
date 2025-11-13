package configs

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email string) (string, error) {
	secret, err := GetEnv("JWT_SECRET")
	if err != nil {
		return "", err
	}

	// ✅ PERBAIKAN: Pakai pointer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenString string) (*JWTClaims, error) {
	secret, err := GetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	// ✅ PERBAIKAN: Validasi signing method
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ✅ PERBAIKAN: Pakai VerifyJWT, jangan duplikat logic
func GetEmailFromToken(tokenString string) (string, error) {
	claims, err := VerifyJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// ✅ TAMBAHKAN: GetIDFromToken yang user minta
func GetIDFromToken(tokenString string) (string, error) {
	claims, err := VerifyJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
