package controllers

import (
	"encoding/json"
	"errors"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/responses"
	"log"
	"net/http"
)

type SignupRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (server *HttpAppServer) Signup(w http.ResponseWriter, r *http.Request) {
	signupRequest := SignupRequest{}
	err := json.NewDecoder(r.Body).Decode(&signupRequest)
	if err != nil {
		log.Printf("Unexpected error while decoding request: %s", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	hash, err := models.HashPassword(signupRequest.Password)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	user := &models.User{
		Login:    signupRequest.Login,
		Password: hash,
	}
	err = models.CreateUser(server.DB, user)

	if err != nil {
		if err, ok := err.(*models.UserError); ok {
			responses.Error(w, http.StatusConflict, err.Err)
		} else {
			responses.Error(w, http.StatusInternalServerError, errors.New("unknown error"))
		}
		return
	}

	session, _ := server.AuthStore.GetSession(r)
	session.Values[auth.SESSION_USER_ID] = user.ID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses.Ok(w)
}
