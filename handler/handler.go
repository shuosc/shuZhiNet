package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/model/student"
	StudentService "shuZhiNet/service/student"
	"shuZhiNet/service/token"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	loginStudent, err := StudentService.Login(input.Username, input.Password)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var output struct {
		StudentName string `json:"student_name"`
		Token       string `json:"token"`
	}
	output.StudentName = loginStudent.Name
	output.Token = token.GenerateJWT(loginStudent)
	student.Save(loginStudent)
	outputJSON, _ := json.Marshal(output)
	_, _ = w.Write(outputJSON)
	log.Println(loginStudent.Name, "logged in.")
}

func ActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("participating") == "true" {
			ParticipatingActivitiesHandler(w, r)
		} else {
			AllActivitiesHandler(w, r)
		}
	case "POST":
		TakePartHandler(w, r)
	case "DELETE":
		OptOutHandler(w, r)
	}
}

func ScholarShipsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ApplyScholarshipHandler(w, r)
	}
}
