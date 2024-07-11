package storage

import (
	"database/sql"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, error, int) {
	const path = "/storage/get_short_forecast"
	var country string
	queryCity := "SELECT country FROM cities where name = $1"
	if err := s.Db.Get(&country, queryCity, city); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to get city from db. This city isn't in forecast list", slog.String("path", path))
			return models.ShortForecastResponse{}, nil, http.StatusBadRequest
		}
		log.Error("Get cities from db error", slog.String("path", path))
		return models.ShortForecastResponse{}, err, http.StatusInternalServerError
	}
	if country == "" {
		log.Error("No prediction for this city")
		return models.ShortForecastResponse{}, nil, http.StatusBadRequest
	}
	log.Debug("Get city, country from db success")

	var weatherList []models.Forecast
	queryWeather := "SELECT temperature, date FROM weather where city_name = $1"
	if err := s.Db.Select(&weatherList, queryWeather, city); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to get forecast from db. No forecast for this city", slog.String("path", path))
			return models.ShortForecastResponse{}, nil, http.StatusBadRequest
		}
		log.Error("Get weather for city from db error", slog.String("path", path))
		return models.ShortForecastResponse{}, err, http.StatusInternalServerError
	}
	log.Debug("Get city, country from db success")
	var meanTemp float64
	var dates []time.Time
	for index := range weatherList {
		meanTemp = meanTemp + weatherList[index].Temp
		dates = append(dates, weatherList[index].Date)
	}
	meanTemp = meanTemp / float64(len(weatherList))
	meanTemp = math.Floor(meanTemp*100) / 100

	//сортировка
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Unix() < dates[j].Unix()
	})
	prediction := models.ShortForecastResponse{city, country, meanTemp, dates}

	return prediction, nil, http.StatusOK
}
