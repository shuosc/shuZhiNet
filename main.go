package main

import (
	"log"
	"net/http"
	"shuZhiNet/handler/activityhandler"
	"shuZhiNet/handler/studenthandler"
)

func main() {
	http.HandleFunc("/login", studenthandler.LoginHandler)
	http.HandleFunc("/activities", activityhandler.ListHandler)
	http.HandleFunc("/engage", activityhandler.EngageHandler)
	http.HandleFunc("/cancel", activityhandler.CancelHandle)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
