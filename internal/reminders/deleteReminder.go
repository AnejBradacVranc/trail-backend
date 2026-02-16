package reminders

import (
	"backend/api"
	"backend/internal/tools"
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func DeleteReminder(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {
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

	if err := (*db).DeleteReminder(reminderID); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	api.RequestSuccessHandler(w, "Reminder deleted")
}
