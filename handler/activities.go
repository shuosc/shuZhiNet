package handler

import "net/http"

func ActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("participating") == "true" {
			ParticipatingActivitiesHandler(w, r)
		} else {
			AllActivitiesHandler(w, r)
		}
	case "POST":
		TakePartHandler(w, r)
	case "DELETE":
		OptOutHandler(w, r)
	}
}

func ScholarShipsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ApplyScholarshipHandler(w, r)
	}
}
