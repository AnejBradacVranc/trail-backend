package handlers

import (
	"backend/api"
	"backend/internal/applicationStatuses"
	"backend/internal/applications"
	"backend/internal/platforms"
	"backend/internal/services"
	"backend/internal/statistics"
	"backend/internal/tools"
	"context"
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

	ctx := context.Background()

	if err != nil {
		log.Error(err)
		return
	}

	s := services.NewTransactionService(ctx, *db)

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		err = json.NewEncoder(w).Encode(api.Response[string]{
			Code: http.StatusOK,
			Data: "Hello I am Backend",
		})
		if err != nil {
			log.Error(err)
		}
	})

	r.HandleFunc("GET /platforms", func(w http.ResponseWriter, r *http.Request) {
		platforms.GetPlatforms(w, r, db)
	})

	r.HandleFunc("GET /user/applications", func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplications(w, r, db)
	})

	r.HandleFunc("GET /user/statistics/summary", func(w http.ResponseWriter, r *http.Request) {
		statistics.GetSummary(w, r, db)
	})

	r.HandleFunc("GET /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplication(w, r, db)
	})

	r.HandleFunc("GET /applications/statuses", func(w http.ResponseWriter, r *http.Request) {
		applicationStatuses.GetStatuses(w, r, db)
	})

	r.HandleFunc("POST /applications", func(w http.ResponseWriter, r *http.Request) {
		applications.CreateApplication(w, r, s)
	})
}
