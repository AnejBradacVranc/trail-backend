package applications

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetApplication(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	applicationId := r.URL.Query().Get("application_id")

	if applicationId == "" {
		log.Error("Application id is empty")
		api.RequestErrorHandler(w, errors.New("application id cannot be empty"))
		return
	}

	applicationIdNum, err := strconv.ParseInt(applicationId, 10, 64)

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	application, err := (*db).GetApplicationByID(applicationIdNum)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(api.Response[tools.ApplicationDetail]{
		Code: http.StatusOK,
		Data: *application,
	})

}
