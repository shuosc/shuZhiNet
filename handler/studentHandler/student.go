package studentHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shuZhiNet/model/student"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/login"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	loginStudent := login.Login(input.Username, input.Password)
	var output struct {
		StudentName string `json:"student_name"`
		Token       string `json:"token"`
	}
	output.StudentName = loginStudent.Name
	output.Token = auth.GenerateJWT(loginStudent)
	student.Save(loginStudent)
	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
}
