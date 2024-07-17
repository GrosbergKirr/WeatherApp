package app_api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// CitiesGetter godoc
// @Summary Get list of available cities
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   page     query    string     true  "Page number"
// @Param   per_page query    string     true  "Number of items per page"
// @Success 200 {object} models.CitiesListResponse "List of cities"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Router /get_cities [get]
func CitiesGetter(log *slog.Logger, cities DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "/app_api/cities_getter"
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Error("Pagination error. Page should be integer", slog.String("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		perPageInt, err := strconv.Atoi(perPage)
		if err != nil {
			log.Error("Pagination error. perPage should be integer", slog.String("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		offset := (pageInt - 1) * perPageInt
		citiesList, stat, err := cities.GetCities(log, perPageInt, offset)
		if err != nil {
			log.Error("Failed to get cities", slog.Any("err", err), slog.String("path", path))
			w.WriteHeader(stat)
			return
		}
		log.Info("Get cities success")
		render.JSON(w, r, citiesList)
	}
}
