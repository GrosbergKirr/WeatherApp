package storage

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) WeatherUpdater(log *slog.Logger, weatherList []models.Forecast) {
	tx, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction", slog.Any("error", err))
		return
	}
	for _, weather := range weatherList {
		query := "UPDATE weather SET temperature = $1, date = $2, full_forecast = $3 WHERE city_id= $4"
		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Error("Could not prepare statement", slog.Any("error", err))
			err = tx.Rollback()
			if err != nil {
				log.Error("Failed rollback", slog.Any("err", err))
				return
			}
			return
		}
		_, err = stmt.Exec(weather.Temp, weather.Date, weather.FullForecast, weather.CityID)
		if err != nil {
			log.Error("Could not update forecast", slog.Any("error", err))
			err = tx.Rollback()
			if err != nil {
				log.Error("Failed rollback", slog.Any("err", err))
				return
			}
			return
		}
		log.Debug("Updated forecast", slog.Any("weather", weather.CityName))
	}
	if err = tx.Commit(); err != nil {
		log.Error("Failed to commit transaction", slog.Any("err", err))
		return
	}
	log.Info("Updated forecast list")
	return
}
