package shuZhiNet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/crawl"
	"shuZhiNet/service/login"
)

func allowOrigin(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	return w
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	student := login.Login(input.Username, input.Password)
	var output struct {
		StudentName string `json:"student_name"`
		Token       string `json:"token"`
	}
	output.StudentName = student.Name
	output.Token = auth.GenerateJWT(student)
	fmt.Println(output.StudentName, "logged in.")
	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
	fmt.Println(crawl.FetchActivitiesByStudent())
}

func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(crawl.FetchActivitiesByStudent())
	w.Write(response)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/activities", activitiesHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
