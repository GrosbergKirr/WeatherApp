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
	cities []models.City) []models.Forecast {
	var weatherList []models.Forecast
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	log.Info("Start parsing cities weather")
	for _, cityIndex := range cities {
		wg.Add(1)
		go func() {
			defer wg.Done()
			statCode, err := GetWeather(log, cli, mu, cityIndex, &weatherList, cfg.WeatherUrl, cfg.ApiKey)
			if err != nil {
				log.Error("Get Forecast fo each city error", slog.String("error", err.Error()),
					slog.Int("error", statCode))
				return
			}
		}()
	}
	log.Info("Parsing cities weather success")
	wg.Wait()
	return weatherList
}
