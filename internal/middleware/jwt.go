package middleware

import (
	"auth-app/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("SECRET_KEY")

func GenerateToken(user *model.User) string {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"roles": user.Roles,
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretKey))

	return signedToken
}

func VerifyToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unauthorized")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.Wrapf(err, "unauthorized")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func ComparePassword(dbPassword, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(userPassword))

	return err == nil
}
