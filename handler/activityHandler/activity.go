package activityHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shuZhiNet/model/activity"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/cancel"
	"shuZhiNet/service/crawl"
	"shuZhiNet/service/engage"
)

func ActivityHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resultObject, _ := activity.Get(id)
	response, _ := json.Marshal(resultObject)
	w.Write(response)
}

func AllActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(crawl.FetchActivitiesByStudent())
	w.Write(response)
}

func ParticipatingActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")[7:]
	student, err := auth.GetStudent(tokenString)
	if err != nil {
		w.WriteHeader(403)
	}
	activityList := crawl.FetchParticipatingActivityIds(student)
	response, _ := json.Marshal(activityList)
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
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	tokenString := r.Header.Get("Authorization")[7:]
	student, err := auth.GetStudent(tokenString)
	if err != nil {
		w.WriteHeader(403)
	}
	participatingActivities := crawl.FetchParticipatingActivityIds(student)
	for _, activityObject := range participatingActivities {
		if activityObject.Id == input.ActivityId {
			cancel.Cancel(student, activityObject.ParticipateInfoId)
		}
	}
}
