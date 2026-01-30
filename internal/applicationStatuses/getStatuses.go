package applicationStatuses

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetStatuses(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	statuses, err := (*db).GetStatuses()

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[[]*tools.ApplicationStatus]{
		Success: true,
		Code:    http.StatusOK,
		Data:    statuses,
	}); err != nil {
		log.Error(err)
	}

}
