package app_client

import (
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func GetLocationApp(log *slog.Logger,
	cfg *internal.Config,
	cli http.Client,
	cityList []string) []models.City {
	log.Info("Start parsing cities info")
	var cities []models.City
	for _, city := range cityList {
		cityResponse, err, code := GetCitiesLocation(log, cli, cfg.CitiesUrl, city, cfg.ApiKey)
		if err != nil {
			log.Error("get location error", slog.Int("Status Code", code), slog.String("error", err.Error()))
			return nil
		}
		cities = append(cities, cityResponse)
	}
	log.Info("Parse cities info success")
	return cities

}
