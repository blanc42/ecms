package utils

import (
	"os"
	"time"

	"github.com/blanc42/ecms/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a new JWT token for the admin
func GenerateToken(admin models.Admin) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin_id"] = admin.ID
	claims["username"] = admin.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		adminID := uint(claims["admin_id"].(float64))
		return adminID, nil
	}

	return 0, jwt.ErrSignatureInvalid
}
