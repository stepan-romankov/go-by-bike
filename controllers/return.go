package controllers

import (
	"errors"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stepan-romankov/go-by-bike/responses"
	"log"
	"net/http"
)

func (server *HttpAppServer) Return(w http.ResponseWriter, r *http.Request) {
	userId, err := server.AuthStore.GetUserId(r)
	if err != nil {
		log.Printf("Unauthorized access")
		responses.Error(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err = models.Return(server.DB, userId)

	if err != nil {
		if err, ok := err.(*models.RentalError); ok {
			responses.Error(w, http.StatusConflict, err.Err)
		} else {
			responses.Error(w, http.StatusInternalServerError, errors.New("unknown error"))
		}
		return
	}

	responses.Ok(w)
}
