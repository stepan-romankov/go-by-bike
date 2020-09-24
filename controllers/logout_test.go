package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	body := bytes.NewBuffer([]byte(`{"login":"user1", "password": "test"}`))
	req, _ := http.NewRequest(http.MethodPost, "/logon", body)
	response := executeRequest(req)
	cookieStr := response.Header()["Set-Cookie"]

	request := &http.Request{Header: http.Header{"Cookie": cookieStr}}

	// Extract the dropped cookie from the request.
	cookie, _ := request.Cookie("auth")

	req, _ = http.NewRequest(http.MethodPost, "/logout", body)
	req.AddCookie(cookie)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
