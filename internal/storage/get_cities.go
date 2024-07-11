package storage

import (
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) GetCities(log *slog.Logger, perPage, offset int) (models.CitiesListResponse, error, int) {
	const path = "/storage/get_cities"
	var cities []models.City
	query := "SELECT name FROM cities LIMIT $1 OFFSET $2"
	if err := s.Db.Select(&cities, query, perPage, offset); err != nil {
		log.Error("Failed to get cities from db ", slog.String("path", path))
		return models.CitiesListResponse{}, err, http.StatusInternalServerError
	}
	var cityList models.CitiesListResponse
	for i := range cities {
		cityList.Cities = append(cityList.Cities, cities[i].Name)
	}
	log.Debug("Get data from db success")
	return cityList, nil, http.StatusOK
}
