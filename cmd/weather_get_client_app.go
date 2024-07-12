package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/app_client"
	"github.com/GrosbergKirr/WeatherApp/internal/storage"
)

func RanWeatherClientApp(ctx context.Context, wg *sync.WaitGroup, log *slog.Logger,
	cfg *internal.Config,
	cli http.Client,
	db *storage.Storage,
	cityList []string) {
	citiesCoordinates := app_client.GetLocationApp(log, cfg, cli, cityList)
	weatherList := app_client.GetWeatherApp(log, cfg, cli, citiesCoordinates)
	db.SaveCitiesToDB(log, citiesCoordinates, weatherList)
	db.SaveWeatherToDB(log, weatherList)

	// Асинхронное обновление погоды
	ctxSig, _ := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	wg.Add(1)
	go func() {
		log.Debug("Start async updating func")
		t := time.Now()
		for {
			select {
			case <-ctxSig.Done():
				log.Info("Updater func gracefully stopped")
				wg.Done()
				return
			default:
				if time.Since(t) > time.Minute {
					log.Info("Starting update")
					weatherList = app_client.GetWeatherApp(log, cfg, cli, citiesCoordinates)
					db.WeatherUpdater(log, weatherList)
					log.Info("Update complete")
					t = time.Now()
				}
			}
		}
	}()
}
