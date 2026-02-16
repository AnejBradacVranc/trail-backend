package reminders

import (
	"backend/api"
	"backend/internal/tools"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type UpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	RemindAt    *string `json:"remind_at"`
	IsCompleted *bool   `json:"is_completed"`
}

func UpdateReminder(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {
	idParam := r.PathValue("id")
	if idParam == "" {
		api.RequestErrorHandler(w, errors.New("reminder id cannot be empty"))
		return
	}
	reminderID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid reminder id"))
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	var parsedRemindAt *time.Time

	if req.RemindAt != nil {
		parsed, err := utils.ParseTimeString(*req.RemindAt)
		parsedRemindAt = &parsed
		if err != nil {
			api.RequestErrorHandler(w, errors.New("invalid remind_at format"))
			return
		}
	}

	if err := (*db).UpdateReminder(reminderID, req.Title, req.Description, parsedRemindAt, req.IsCompleted); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	api.RequestSuccessHandler(w, "Reminder updated")
}
