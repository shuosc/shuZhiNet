package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/model/student"
	"shuZhiNet/service/activity"
	"shuZhiNet/service/token"
)

func AllActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	activityList := activity.FetchActivities()
	// Go is kidding
	// @see https://github.com/golang/go/issues/26866
	if len(activityList) == 0 {
		_, _ = w.Write([]byte("[]"))
	} else {
		response, _ := json.Marshal(activityList)
		_, _ = w.Write(response)
	}
}

func getStudent(r *http.Request) (student.Student, error) {
	tokenInHeader := r.Header.Get("Authorization")
	if tokenInHeader == "" {
		return student.Student{}, errors.New("no token given")
	}
	tokenString := tokenInHeader[7:]
	return token.GetStudent(tokenString)
}

func ParticipatingActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	activityList := activity.FetchParticipatingActivities(studentObject)
	// Go is kidding
	// @see https://github.com/golang/go/issues/26866
	if len(activityList) == 0 {
		_, _ = w.Write([]byte("[]"))
	} else {
		response, _ := json.Marshal(activityList)
		_, _ = w.Write(response)
	}
}

func TakePartHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var input struct {
		ActivityId  string `json:"id"`
		PhoneNumber string `json:"phone_number"`
		MailAddress string `json:"mail_address"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	activity.TakePartInActivity(studentObject, input.ActivityId, input.PhoneNumber, input.MailAddress)
	log.Println(studentObject.Name, "take part in activity", input.ActivityId)
}

func OptOutHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	id := r.URL.Query().Get("id")
	participatingActivities := activity.FetchParticipatingActivities(studentObject)
	for _, activityObject := range participatingActivities {
		if activityObject.Id == id {
			activity.OptOutActivity(studentObject, activityObject.ParticipateInfoId)
			log.Println(studentObject.Name, "opt out of activity", activityObject.Id)
			break
		}
	}
}
