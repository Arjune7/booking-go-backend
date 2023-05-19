package helper

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// type authHelp interface {
// 	HandleHashPassword(password string) (string, error)
// 	HandleComparePassword(password, hashedPassword string) bool
// }

func HandleHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func HandleComparePassword(password, hashedPassword string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
	fmt.Println(res)
	return res == nil
}

func HandleGenerateToken(email, password string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secretKey := []byte("secret124")

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token : %v", err)
	}
	return signedToken, nil
}
