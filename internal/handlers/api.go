package handlers

import (
	"backend/api"
	"backend/internal/applicationStatuses"
	"backend/internal/applications"
	"backend/internal/auth"
	"backend/internal/platforms"
	"backend/internal/services"
	"backend/internal/statistics"
	"backend/internal/tools"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			api.RequestUnauthorisedHandler(w, "Unauthorized")
			return
		}

		token, err := jwt.ParseWithClaims(
			cookie.Value,
			&jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			api.RequestUnauthorisedHandler(w, "Invalid token")
			return
		}

		claims := token.Claims.(*jwt.RegisteredClaims)
		ctx := context.WithValue(r.Context(), "user_id", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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

	r.HandleFunc("POST /user/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, db)
	})

	r.HandleFunc("POST /user/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(w, r, db)
	})

	r.HandleFunc("POST /user/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(w, r)
	})
	r.Handle("GET /user/check", Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.Check(w, r, db)
	})))

	r.Handle("GET /user/applications", Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplications(w, r, db)
	})))

	r.Handle("GET /user/statistics/summary", Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statistics.GetSummary(w, r, db)
	})))

	r.Handle("GET /applications/{id}", Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applications.GetApplication(w, r, db)
	})))

	r.HandleFunc("GET /applications/statuses", func(w http.ResponseWriter, r *http.Request) {
		applicationStatuses.GetStatuses(w, r, db)
	})

	r.Handle("POST /applications", Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applications.CreateApplication(w, r, s)
	})))

}
