package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"shuZhiNet/model/student"
)

func GenerateJWT(student student.Student) string {
	fmt.Println("generating jwt")
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": student.Id,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	fmt.Println(result)
	fmt.Println("generate jwt end")
	return result
}
