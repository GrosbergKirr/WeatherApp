package app_client

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/GrosbergKirr/WeatherApp/internal/models"

	"github.com/go-chi/render"
)

func GetCitiesLocation(log *slog.Logger, client http.Client, sideApiUrl string, cityName string, apiKey string) (models.City, error, int) {
	u, err := url.Parse(sideApiUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return models.City{}, err, http.StatusInternalServerError
	}
	q := u.Query()
	q.Set("q", cityName)
	q.Set("limit", "5")
	q.Set("appid", apiKey)
	u.RawQuery = q.Encode()
	urlBody := u.String()

	resp, err := client.Get(urlBody)
	if err != nil {
		log.Error("failed to get response")
		return models.City{}, err, resp.StatusCode

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
