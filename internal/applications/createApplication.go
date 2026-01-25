package applications

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// CreateApplication handles POST /applications
// Expected JSON body:
//
//	{
//	  "user_id": 1,
//	  "status_id": 2,
//	  "company_id": 3,
//	  "job_title": "Software Engineer",
//	  "platform": "LinkedIn",
//	  "job_url": "https://example.com/job",
//	  "salary_min": 60000,
//	  "salary_max": 90000,
//	  "applied_at": "2026-01-25T19:23:00Z" // RFC3339
//	}

type Request struct {
	UserID      int64   `json:"user_id"`
	StatusID    int64   `json:"status_id"`
	CompanyName string  `json:"company_name"`
	Location    *string `json:"location"`
	JobTitle    string  `json:"job_title"`
	Platform    string  `json:"platform"`
	JobURL      *string `json:"job_url"`
	SalaryMin   *int    `json:"salary_min"`
	SalaryMax   *int    `json:"salary_max"`
	AppliedAt   string  `json:"applied_at"`
}

func CreateApplication(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	//TODO Create company
	//TODO Insert notes

	if r.Method != http.MethodPost {
		api.RequestErrorHandler(w, errors.New("invalid method"))
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	if req.UserID == 0 || req.StatusID == 0 || req.CompanyName == "" || req.JobTitle == "" || req.Platform == "" || req.AppliedAt == "" {
		api.RequestErrorHandler(w, errors.New("missing required fields"))
		return
	}

	appliedAt, err := time.Parse(
		"2006-01-02 15:04:05.999999-07",
		req.AppliedAt,
	)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("applied_at must be RFC3339 timestamp"))
		return
	}

	//TODO CREATE COMPANY AND INSERT ID HERE, THE SAME FOR NOTES, THEY ARE IN THE OTHER TABLE
	id, err := (*db).CreateApplication(req.UserID, req.StatusID, 2, req.JobTitle, req.Platform, req.JobURL, req.SalaryMin, req.SalaryMax, appliedAt)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(api.Response[api.IDResp]{
		Code: http.StatusOK,
		Data: api.IDResp{ID: id},
	}); err != nil {
		log.Error(err)
	}
}
