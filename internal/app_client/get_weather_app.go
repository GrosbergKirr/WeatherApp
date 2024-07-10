package app_client

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func GetWeatherApp(log *slog.Logger,
	cfg *internal.Config,
	cli http.Client,
	cities []models.City) []models.Weather {
	var respStruct models.WeatherResponse
	var weatherList []models.Weather
	wg := new(sync.WaitGroup)
	log.Info("Start parsing cities weather")
	for cityIndex := range cities {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err, statCode := GetWeather(log, cli, respStruct, cities[cityIndex], &weatherList, cfg.WeatherUrl, cfg.ApiKey)
			if err != nil {
				log.Error("Get Weather fo each city error", slog.String("error", err.Error()),
					slog.Int("error", statCode))
				return
			}
		}()
	}
	log.Info("Parsing cities weather success")
	wg.Wait()
	return weatherList
}
