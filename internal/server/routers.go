package server

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/app_api"
	"github.com/GrosbergKirr/WeatherApp/internal/storage"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRouters(log *slog.Logger, db *storage.Storage) chi.Router {
	router := chi.NewRouter()
	router.Get("/get_cities", app_api.CitiesGetter(log, db))
	router.Get("/get_short_forecast", app_api.ShortPredGetter(log, db))
	router.Post("/get_full_forecast", app_api.FullPredGetter(log, db))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9090/swagger/doc.json"), //The url pointing to API definition
	))
	log.Info("Set routers successfully")
	log.Info("USE WEB SWAGGER ON: http://localhost:9090/swagger/index.html#/ ")
	return router
}
