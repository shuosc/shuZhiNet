package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shuZhiNet/service/scholarship"
)

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
