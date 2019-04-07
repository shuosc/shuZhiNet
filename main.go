package main

import (
	"log"
	"net/http"
	"os"
	"shuZhiNet/handler"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/all-activities", handler.AllActivitiesHandler)
	http.HandleFunc("/participating-activities", handler.ParticipatingActivitiesHandler)
	http.HandleFunc("/take-part", handler.TakePartHandler)
	http.HandleFunc("/opt-out", handler.OptOutHandler)
	http.HandleFunc("/apply-scholarship", handler.ApplyScholarshipHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
