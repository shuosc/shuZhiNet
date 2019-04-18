package main

import (
	"log"
	"net/http"
	"os"
	"shuZhiNet/handler"
)

func main() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/ping", handler.PingPongHandler)
	http.HandleFunc("/activities", handler.ActivitiesHandler)
	http.HandleFunc("/scholarships", handler.ScholarShipsHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
