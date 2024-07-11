package storage

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) SaveCitiesToDB(log *slog.Logger, cities []models.City, weatherList []models.Forecast) {
	tx1, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction", slog.Any("err", err))
		return
	}
	queryCities := "INSERT INTO cities (name, country, latitude, longitude) VALUES(:name, :country, :latitude, :longitude) RETURNING id"
	for i := range cities {
		stmt, err := tx1.PrepareNamed(queryCities)
		if err != nil {
			log.Error("Failed to prepare cities transactions to DB", slog.Any("err", err))
			return
		}
		defer func(stmt *sqlx.NamedStmt) {
			err := stmt.Close()
			if err != nil {
				log.Error("failed to close statement", slog.Any("err", err))
			}
		}(stmt)

		err = stmt.QueryRowx(cities[i]).Scan(&cities[i].Id)
		if err != nil {
			log.Error("Failed to insert into DB", slog.Any("err", err))
			return
		}
		log.Debug("Saved a city to DB", slog.String("city", cities[i].Name))
	}
	for i := range weatherList {
		for j := range cities {
			if weatherList[i].CityName == cities[j].Name {
				weatherList[i].CityID = cities[j].Id
			}
		}
	}
	log.Info("Save cities to DB successfully")
	if err = tx1.Commit(); err != nil {
		log.Error("Failed to commit transaction", slog.Any("err", err))
		return
	}
}
func (s *Storage) SaveWeatherToDB(log *slog.Logger, weatherList []models.Forecast) {
	tx2, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction", slog.Any("err", err))
		return
	}
	queryWeather := "INSERT INTO weather (city_name, temperature, date, city_id, full_forecast) VALUES(:city_name, :temperature, :date, :city_id, :full_forecast)"
	for _, weather := range weatherList {
		stmt, err := tx2.PrepareNamed(queryWeather)
		if err != nil {
			log.Error("Failed to prepare weather transactions to DB", slog.Any("err", err))
			err := tx2.Rollback()
			if err != nil {
				log.Error("Failed rollback", slog.Any("err", err))
				return
			}
			return
		}
		defer func(stmt *sqlx.NamedStmt) {
			err := stmt.Close()
			if err != nil {
				log.Error("failed to close statement", slog.Any("err", err))
				return
			}
		}(stmt)

		_, err = stmt.Exec(weather)
		if err != nil {
			log.Error("Failed to insert into DB", slog.Any("err", err))
			err := tx2.Rollback()
			if err != nil {
				log.Error("Failed rollback", slog.Any("err", err))
				return
			}
			return
		}
		log.Debug("Success save a weather to DB", slog.String("city", weather.CityName), slog.String("date", weather.Date.String()))
	}
	if err = tx2.Commit(); err != nil {
		log.Error("Failed to commit transaction", slog.Any("err", err))
		return
	}
	log.Info("Saved data to Db success")
	return
}
