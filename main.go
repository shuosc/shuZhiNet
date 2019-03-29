package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/model/student"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/crawl"
	"shuZhiNet/service/engage"
	"shuZhiNet/service/login"
)

func allowOrigin(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	return w
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
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
	fmt.Println(output.StudentName, "logged in.")
	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
	fmt.Println(crawl.FetchActivitiesByStudent())
}

func engageHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
	var input struct {
		ActivityId  string `json:"activity_id"`
		PhoneNumber string `json:"phone_number"`
		MailAddress string `json:"mail_address"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	tokenString := r.Header.Get("Authorization")[7:]
	student, err := auth.GetStudent(tokenString)
	fmt.Println(student)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	engage.Engage(student, input.ActivityId, input.PhoneNumber, input.MailAddress)
}

func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
	response, _ := json.Marshal(crawl.FetchActivitiesByStudent())
	w.Write(response)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/activities", activitiesHandler)
	http.HandleFunc("/engage", engageHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
