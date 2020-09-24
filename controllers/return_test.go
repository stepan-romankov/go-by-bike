package controllers

import (
	"bytes"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturn(t *testing.T) {
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
	cookie1 := login("user1", "test")

	response := rentBike(cookie1)
	checkResponseCode(t, http.StatusOK, response.Code)

	response = returnBike(cookie1)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestReturnNotRented(t *testing.T) {
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
	cookie1 := login("user1", "test")

	response := returnBike(cookie1)
	checkResponseCode(t, http.StatusConflict, response.Code)
}

func TestReturnAlreadyReturned(t *testing.T) {
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
	cookie1 := login("user1", "test")
	_ = rentBike(cookie1)
	_ = returnBike(cookie1)
	response := returnBike(cookie1)
	checkResponseCode(t, http.StatusConflict, response.Code)
}

func rentBike(cookie1 *http.Cookie) *httptest.ResponseRecorder {
	body := bytes.NewBuffer([]byte(`{"bike_id": 1}`))
	req, _ := http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie1)
	response := executeRequest(req)
	return response
}

func returnBike(cookie1 *http.Cookie) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/return", nil)
	req.AddCookie(cookie1)
	response := executeRequest(req)
	return response
}
