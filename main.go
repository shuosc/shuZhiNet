package main

import (
	"log"
	"net/http"
	"shuZhiNet/handler/activityHandler"
	"shuZhiNet/handler/studentHandler"
)

func main() {
	http.HandleFunc("/login", studentHandler.LoginHandler)
	http.HandleFunc("/all-activities", activityHandler.AllActivitiesHandler)
	http.HandleFunc("/participating-activities", activityHandler.ParticipatingActivitiesHandler)
	http.HandleFunc("/activity", activityHandler.ActivityHandler)
	http.HandleFunc("/engage", activityHandler.EngageHandler)
	http.HandleFunc("/cancel", activityHandler.CancelHandle)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
