package controllers

import (
	"encoding/json"
	"github.com/stepan-romankov/go-by-bike/api/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestBikes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/bikes", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var bikes []models.Bike
	err := json.NewDecoder(response.Body).Decode(&bikes)
	if err != nil {
		panic(err)
	}

	require.Len(t, bikes, 3)
}

func TestBikeIsRented(t *testing.T) {
	Server.DB.Model(&models.Rental{UserID: 1, BikeId: 1}).Insert()

	req, _ := http.NewRequest(http.MethodGet, "/bikes", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var bikes []models.Bike
	err := json.NewDecoder(response.Body).Decode(&bikes)
	if err != nil {
		panic(err)
	}

	for _, v := range bikes {
		if v.ID == 1 {
			require.True(t, v.Rented)
		}
	}
	Server.DB.Model(&models.Rental{}).Where("1=1").Delete()
}
