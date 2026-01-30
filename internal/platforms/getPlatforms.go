package platforms

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetPlatforms(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	platforms, err := (*db).GetPlatforms()

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[[]*tools.Platform]{
		Success: true,
		Code:    http.StatusOK,
		Data:    platforms,
	}); err != nil {
		log.Error(err)
	}

}
