package main

import (
	"flag"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/cmd"
	_ "github.com/GrosbergKirr/WeatherApp/docs"
	"github.com/GrosbergKirr/WeatherApp/internal"
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

var configPath string
var cityList = []string{
	"Moscow", "Paris", "Berlin", "London",
	"Madrid", "Rome", "Washington", "Ottawa",
	"Minsk", "Tokyo", "Canberra", "Tallinn",
	"Warsaw", "Budapest", "Jakarta", "Prague",
	"Lisbon", "Beijing", "Ankara", "Seoul",
}

func init() {
	flag.StringVar(&configPath, "c", "./config/config.yaml", "Path to config file")
}
func main() {
	flag.Parse()
	log := internal.SetupLogger()
	cfg := internal.SetupConfig(log, configPath)
	db := storage.InitStorage(log, cfg.DBUsername, cfg.DBPassword, cfg.DBAddress, cfg.DBName, cfg.DBMode)
	cli := http.Client{}
	//ctx := context.Background()

	cmd.RanWeatherClientApp(log, cfg, cli, db, cityList)
	//cmd.WeatherServiceApp(ctx, log, cfg, db)

}
