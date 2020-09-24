package controllers

import (
	"encoding/json"
	"errors"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stepan-romankov/go-by-bike/auth"
	"github.com/stepan-romankov/go-by-bike/responses"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type LogonRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogonResponse struct {
	UserID uint32 `json:"user_id"`
}

func (server *HttpAppServer) Logon(w http.ResponseWriter, r *http.Request) {
	logonRequest := LogonRequest{}
	var err = json.NewDecoder(r.Body).Decode(&logonRequest)
	if err != nil {
		log.Printf("Unexpected error while decoding request: %s", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := models.FindUserByLogin(server.DB, logonRequest.Login)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, errors.New("unknown failure"))
		return
	}

	if user == nil {
		responses.Error(w, http.StatusForbidden, errors.New("user not found"))
		return
	}

	err = models.VerifyPassword(user.Password, logonRequest.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		responses.Error(w, http.StatusForbidden, errors.New("invalid password"))
		return
	}

	session, _ := server.AuthStore.GetSession(r)
	session.Values[auth.SESSION_USER_ID] = user.ID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.Json(w, LogonResponse{UserID: user.ID})
}
