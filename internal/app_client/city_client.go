package app_client

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal/models"

	"github.com/go-chi/render"
)

func GetCitiesLocation(log *slog.Logger, client http.Client, sideApiUrl string, cityName string, apiKey string) (models.City, error, int) {
	urlBody := sideApiUrl + fmt.Sprintf("q=%s&limit=5&appid=", cityName) + apiKey

	resp, err := client.Get(urlBody)
	if err != nil {
		log.Error("failed to get response")
		return models.City{}, err, http.StatusInternalServerError

	}
	defer resp.Body.Close()

	var city []models.City
	err = render.DecodeJSON(resp.Body, &city)
	if err != nil {
		log.Error("fail to decode json")
		return models.City{}, err, resp.StatusCode
	}

	return city[0], nil, resp.StatusCode
}
