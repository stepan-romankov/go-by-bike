package controllers

import (
	"encoding/json"
	"errors"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stepan-romankov/go-by-bike/responses"
	"log"
	"net/http"
)

type RentalRequest struct {
	BikeID uint32 `json:"bike_id"`
}

func (server *HttpAppServer) Rent(w http.ResponseWriter, r *http.Request) {
	rentalRequest := RentalRequest{}
	err := json.NewDecoder(r.Body).Decode(&rentalRequest)
	if err != nil {
		log.Printf("Unexpected error while decoding request: %s", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	userID, err := server.AuthStore.GetUserId(r)
	if err != nil {
		responses.Error(w, http.StatusForbidden, err)
		return
	}

	_, err = models.Rent(server.DB, userID, rentalRequest.BikeID)

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
