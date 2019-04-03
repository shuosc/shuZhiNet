package activityhandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/crawl"
	"shuZhiNet/service/engage"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Visited")
	response, _ := json.Marshal(crawl.FetchActivitiesByStudent())
	w.Write(response)
}

func MylistHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")[7:]
	fmt.Println(tokenString)
	student, err := auth.GetStudent(tokenString)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(403)
	}
	fmt.Println(student)
	activityIdList, s:= crawl.FetchMyActivity(student)
	fmt.Println(s,activityIdList)
	response, _ :=json.Marshal(activityIdList)
	w.Write(response)
}

func EngageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ActivityId  string `json:"activity_id"`
		PhoneNumber string `json:"phone_number"`
		MailAddress string `json:"mail_address"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	tokenString := r.Header.Get("Authorization")[7:]
	student, err := auth.GetStudent(tokenString)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	engage.Engage(student, input.ActivityId, input.PhoneNumber, input.MailAddress)
}

func CancelHandle(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ActivityId string `json:"activity_id"`
	}
	var cancelId string
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	tokenString := r.Header.Get("Authorization")[7:]
	student, err := auth.GetStudent(tokenString)
	if err != nil {
		w.WriteHeader(403)
	}
	activityIdList, cancelIdList := crawl.FetchMyActivity(student)
	for i := 0; i < len(activityIdList); i++ {
		if activityIdList[i] == input.ActivityId {
			cancelId = cancelIdList[i]
		}
	}
	engage.Cancel(student, cancelId)
}
