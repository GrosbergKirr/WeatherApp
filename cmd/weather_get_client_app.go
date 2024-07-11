package cmd

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/app_client"
	"github.com/GrosbergKirr/WeatherApp/internal/storage"
)

func RanWeatherClientApp(log *slog.Logger,
	cfg *internal.Config,
	cli http.Client,
	db *storage.Storage,
	cityList []string) {
	citiesCoordinates := app_client.GetLocationApp(log, cfg, cli, cityList)
	weatherList := app_client.GetWeatherApp(log, cfg, cli, citiesCoordinates)
	db.SaveCitiesToDB(log, citiesCoordinates, weatherList)
	db.SaveWeatherToDB(log, weatherList)

	// Асинхронное обновление погоды
	go func() {
		time.Sleep(time.Minute)
		for {
			log.Info("Starting update")
			weatherList = app_client.GetWeatherApp(log, cfg, cli, citiesCoordinates)
			db.WeatherUpdater(log, weatherList)
			log.Info("Update complete")
			time.Sleep(time.Minute)
		}
	}()
}
