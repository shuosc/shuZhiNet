package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func AuthToken(id string, tokenString string) bool {
	fmt.Println("try auth")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	fmt.Println("parsing token")
	if token.Valid {
		fmt.Println("token right")
		return true
	}
	fmt.Println("token left")
	return false
}
