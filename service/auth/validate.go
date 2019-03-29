package auth

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"shuZhiNet/model/student"
)

func GetStudent(tokenString string) (student.Student, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return student.Student{}, nil
	}
	claims := token.Claims.(jwt.MapClaims)
	studentId := claims["studentId"].(string)
	return student.Get(studentId)
}
