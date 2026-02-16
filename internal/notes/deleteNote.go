package notes

import (
	"backend/api"
	"backend/internal/tools"
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func DeleteNote(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {
	idParam := r.PathValue("id")
	if idParam == "" {
		api.RequestErrorHandler(w, errors.New("note id cannot be empty"))
		return
	}
	noteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid note id"))
		return
	}

	if err := (*db).DeleteNote(noteID); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	api.RequestSuccessHandler(w, "Note deleted")
}
