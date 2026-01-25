package handlers

import (
	"backend/api"
	"backend/internal/applications"
	"backend/internal/tools"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Handler(r *http.ServeMux) {

	db, err := tools.NewDatabase()

	if err != nil {
		log.Error(err)
		return
	}

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		err = json.NewEncoder(w).Encode(api.Response[string]{
			Code: http.StatusOK,
			Data: "Hello I am Backend",
		})
		if err != nil {
			log.Error(err)
		}
	})

	r.HandleFunc("GET /applications", func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplications(w, r, db)
	})

	r.HandleFunc("GET /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplication(w, r, db)
	})
}
