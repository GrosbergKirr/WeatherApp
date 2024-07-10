package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GrosbergKirr/WeatherApp/docs"
	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/app_client"
	"github.com/GrosbergKirr/WeatherApp/internal/server"
	"github.com/GrosbergKirr/WeatherApp/internal/storage"
)

// @title Weather app swagger
// @version 1.0
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090

func main() {
	log := internal.SetupLogger()
	cfg := internal.SetupConfig(log)
	db := storage.InitStorage(log, cfg.DBUsername, cfg.DBPassword, cfg.DBAddress, cfg.DBName, cfg.DBMode)
	fmt.Println(cfg.DBAddress)
	cli := http.Client{}
	cityList := []string{
		"Moscow",
		"Paris",
		"Berlin",
		"London",
		"Madrid",
		"Rome",
		"Washington",
		"Ottawa",
		"Minsk",
		"Tokyo",
		"Canberra",
		"Tallinn",
		"Warsaw",
		"Budapest",
		"Jakarta",
		"Prague",
		"Lisbon",
		"Beijing",
		"Ankara",
		"Seoul",
	}

	citiesCoordinates := app_client.GetLocationApp(log, cfg, cli, cityList)
	_ = citiesCoordinates

	weatherList := app_client.GetWeatherApp(log, cfg, cli, citiesCoordinates)

	_ = weatherList
	err := db.SaveToDB(log, citiesCoordinates, weatherList)
	if err != nil {
		log.Error("Error saving cities to database", slog.Any("err", err))
	}
	ctx := context.Background()
	router := server.SetRouters(ctx, log, db)
	newServer := server.NewServer(cfg, router)

	serverStopSig := make(chan os.Signal)
	signal.Notify(serverStopSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go newServer.ServerRun(log, cfg)
	<-serverStopSig
	newServer.ServerStop(ctx, log)

}
