package applications

import (
	"backend/api"
	"backend/internal/services"
	"backend/internal/utils"
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
	UserID        int64   `json:"user_id"`
	StatusID      int64   `json:"status_id"`
	CompanyName   string  `json:"company_name"`
	Location      string  `json:"location"`
	JobTitle      string  `json:"job_title"`
	PlatformID    int64   `json:"platform_id"`
	JobURL        *string `json:"job_url"`
	EmployeeCount *int    `json:"employee_count"`
	SalaryMin     *int    `json:"salary_min"`
	SalaryMax     *int    `json:"salary_max"`
	AppliedAt     string  `json:"applied_at"`
	InterviewAt   *string `json:"interview_at"`
	NoteContent   *string `json:"note_content"`
}

func CreateApplication(w http.ResponseWriter, r *http.Request, s *services.AppServiceInterface) {

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

	if req.UserID == 0 || req.StatusID == 0 || req.CompanyName == "" || req.JobTitle == "" || req.AppliedAt == "" || req.Location == "" {
		api.RequestErrorHandler(w, errors.New("missing required fields"))
		return
	}

	appliedAt, err := utils.ParseTimeString(req.AppliedAt)
	if err != nil {
		api.RequestErrorHandler(w, errors.New("invalid applied_at format"))
		return
	}

	var interviewAt *time.Time

	if req.InterviewAt != nil {
		parsedInterviewAt, err := utils.ParseTimeString(*req.InterviewAt)

		if err != nil {
			api.RequestErrorHandler(w, errors.New("invalid interview_at format"))
			return
		}

		interviewAt = &parsedInterviewAt
	}

	id, err := (*s).CreateApplication(
		req.UserID,
		req.StatusID,
		req.CompanyName,
		req.Location,
		req.EmployeeCount,
		req.JobTitle,
		req.PlatformID,
		req.JobURL,
		req.SalaryMin,
		req.SalaryMax,
		appliedAt,
		interviewAt,
		req.NoteContent)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[api.IDResp]{
		Success: true,
		Code:    http.StatusOK,
		Data:    api.IDResp{ID: *id},
	}); err != nil {
		log.Error(err)
	}
}
