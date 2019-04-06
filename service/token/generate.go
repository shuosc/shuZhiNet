package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"shuZhiNet/model/student"
)

func GenerateJWT(student student.Student) string {
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": student.Id,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return result
}
