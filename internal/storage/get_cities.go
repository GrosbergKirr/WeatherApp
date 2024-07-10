package storage

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) GetCities(log *slog.Logger, page, perPage string) (models.CitiesListResponse, error, int) {
	const path = "/storage/get_cities"
	var cities []models.City

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Error("Pagination error. Page should be integer", slog.String("path", path))
		return models.CitiesListResponse{}, err, http.StatusBadRequest
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		log.Error("Pagination error. perPage should be integer", slog.String("path", path))
		return models.CitiesListResponse{}, err, http.StatusBadRequest
	}
	offset := (pageInt - 1) * perPageInt
	query := "SELECT name FROM cities LIMIT $1 OFFSET $2"
	if err := s.Db.Select(&cities, query, perPageInt, offset); err != nil {
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
