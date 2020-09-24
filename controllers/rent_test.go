package controllers

import (
	"bytes"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"net/http"
	"testing"
)

func TestRentalSuccess(t *testing.T) {
	cookie := login("user1", "test")

	body := bytes.NewBuffer([]byte(`{"bike_id": 1}`))
	req, _ := http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
}

func TestSecondRentalTheSameBikeNotAllowed(t *testing.T) {
	cookie1 := login("user1", "test")
	cookie2 := login("user2", "test")

	body := bytes.NewBuffer([]byte(`{"bike_id": 2}`))
	req, _ := http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie1)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body = bytes.NewBuffer([]byte(`{"bike_id": 2}`))
	req, _ = http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie2)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusConflict, response.Code)
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
}

func TestSecondRentalByTheSameUserNotAllowed(t *testing.T) {
	cookie := login("user1", "test")

	body := bytes.NewBuffer([]byte(`{"bike_id": 1}`))
	req, _ := http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	body = bytes.NewBuffer([]byte(`{"bike_id": 2}`))
	req, _ = http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusConflict, response.Code)
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
}

func TestRentalAfterReturn(t *testing.T) {
	Server.DB.Model(&models.Rental{UserID: 1, BikeId: 1, Completed: true}).Insert()
	cookie := login("user1", "test")

	body := bytes.NewBuffer([]byte(`{"bike_id": 1}`))
	req, _ := http.NewRequest(http.MethodPost, "/rent", body)
	req.AddCookie(cookie)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
