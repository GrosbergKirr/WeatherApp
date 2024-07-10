package storage

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/jmoiron/sqlx"
)

func (s *Storage) SaveToDB(log *slog.Logger, cities []models.City, weatherlist []models.Weather) error {
	tx1, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction")
		return err
	}
	queryCities := "INSERT INTO cities (name, country, latitude, longitude) VALUES(:name, :country, :latitude, :longitude) RETURNING id"
	for city := range cities {
		stmt, err := tx1.PrepareNamed(queryCities)
		if err != nil {
			log.Error("Failed to prepare cities transactions to DB")
			return err
		}
		defer func(stmt *sqlx.NamedStmt) {
			err := stmt.Close()
			if err != nil {
				log.Error("failed to close statement")
			}
		}(stmt)

		err = stmt.QueryRowx(cities[city]).Scan(&cities[city].Id)
		if err != nil {
			log.Error("Failed to insert into DB")
		}
		log.Debug("Saved a city to DB", slog.String("city", cities[city].Name))
	}
	for weather := range weatherlist {
		for _, city := range cities {
			if weatherlist[weather].CityName == city.Name {
				weatherlist[weather].CityID = city.Id
			}
		}
	}
	log.Info("Save cities to DB successfully")
	if err = tx1.Commit(); err != nil {
		log.Error("Failed to commit transaction")
	}

	tx2, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction")
		return err
	}
	queryWeather := "INSERT INTO weather (city_name, temperature, date, city_id, full_forecast) VALUES(:city_name, :temperature, :date, :city_id, :full_forecast)"
	for _, weather := range weatherlist {
		stmt, err := tx2.PrepareNamed(queryWeather)

		if err != nil {
			log.Error("Failed to prepare weather transactions to DB")
			err := tx2.Rollback()
			if err != nil {
				log.Error("Failed rollback")
				return err
			}
			return err
		}
		defer func(stmt *sqlx.NamedStmt) {
			err := stmt.Close()
			if err != nil {
				log.Error("failed to close statement")
			}
		}(stmt)

		_, err = stmt.Exec(weather)
		if err != nil {
			log.Error("Failed to insert into DB")
			err := tx2.Rollback()
			if err != nil {
				log.Error("Failed rollback")
				return err
			}
			return err
		}
		log.Debug("Success save a weather to DB", slog.String("city", weather.CityName), slog.String("date", weather.Date.String()))
	}
	if err = tx2.Commit(); err != nil {
		log.Error("Failed to commit transaction")
		return err
	}
	log.Info("Saved data to Db success")
	return nil
}
