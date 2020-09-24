package controllers

import (
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stepan-romankov/go-by-bike/responses"
	"net/http"
)

func (server *HttpAppServer) Bikes(w http.ResponseWriter, r *http.Request) {
	bikes, err := models.GetBikes(server.DB)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.Json(w, bikes)
}
