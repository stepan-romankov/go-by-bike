package controllers

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	body := bytes.NewBuffer([]byte(`{"login":"ux3", "password": "test"}`))
	req, _ := http.NewRequest(http.MethodPost, "/signup", body)
	_ = executeRequest(req)

	body = bytes.NewBuffer([]byte(`{"login":"ux3", "password": "test"}`))
	req, _ = http.NewRequest(http.MethodPost, "/logon", body)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	cookie := response.Header().Get("Set-Cookie")
	require.NotEmpty(t, cookie)
}

func TestInvalidCredentials(t *testing.T) {
	body := bytes.NewBuffer([]byte(`{"login":"___", "password": "___"}`))
	req, _ := http.NewRequest(http.MethodPost, "/logon", body)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusForbidden, response.Code)
	cookie := response.Header().Get("Set-Cookie")
	require.Empty(t, cookie)
}
