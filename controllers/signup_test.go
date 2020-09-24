package controllers

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSignupSuccess(t *testing.T) {
	var body = bytes.NewBuffer([]byte(`{"login":"ux1", "password": "test"}`))
	req, _ := http.NewRequest(http.MethodPost, "/signup", body)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	cookie := response.Header().Get("Set-Cookie")
	require.NotEmpty(t, cookie)
}

func TestSignupDuplicate(t *testing.T) {
	var body = bytes.NewBuffer([]byte(`{"login":"ux2", "password": "test"}`))
	req, _ := http.NewRequest(http.MethodPost, "/signup", body)
	_ = executeRequest(req)

	body = bytes.NewBuffer([]byte(`{"login":"ux2", "password": "test"}`))
	req, _ = http.NewRequest(http.MethodPost, "/signup", body)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusConflict, response.Code)
	cookie := response.Header().Get("Set-Cookie")
	require.Empty(t, cookie)
}
