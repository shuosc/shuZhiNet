package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/model/student"
	ActivityService "shuZhiNet/service/activity"
	"shuZhiNet/service/scholarship"
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

func AllActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(ActivityService.FetchActivities())
	_, _ = w.Write(response)
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
	activityList := ActivityService.FetchParticipatingActivities(studentObject)
	response, _ := json.Marshal(activityList)
	_, _ = w.Write(response)
}

func TakePartHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var input struct {
		ActivityId  string `json:"activity_id"`
		PhoneNumber string `json:"phone_number"`
		MailAddress string `json:"mail_address"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	ActivityService.TakePartInActivity(studentObject, input.ActivityId, input.PhoneNumber, input.MailAddress)
	log.Println(studentObject.Name, "take part in activity", input.ActivityId)
}

func OptOutHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var input struct {
		ActivityId string `json:"activity_id"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	participatingActivities := ActivityService.FetchParticipatingActivities(studentObject)
	for _, activityObject := range participatingActivities {
		if activityObject.Id == input.ActivityId {
			ActivityService.OptOutActivity(studentObject, activityObject.ParticipateInfoId)
			log.Println(studentObject.Name, "opt out of activity", activityObject.Id)
			break
		}
	}
}

func ApplyScholarshipHandler(w http.ResponseWriter, r *http.Request) {
	studentObject, err := getStudent(r)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var input struct {
		ScholarshipType int    `json:"scholarship_type"`
		Qualifications  []int  `json:"qualifications"`
		Reason          string `json:"reason"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)
	scholarship.Apply(studentObject, input.ScholarshipType, input.Qualifications, input.Reason)
}
