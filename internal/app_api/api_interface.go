package app_api

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

type WeatherInterface interface {
	GetCities(log *slog.Logger, page, perPage string) (models.CitiesListResponse, error, int)
	GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, error, int)
	GetFullPred(log *slog.Logger, param models.Weather) (models.Record, error, int)
}
