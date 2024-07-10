package app_api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// ShortPredGetter godoc
// @Summary Get list of available cities
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   city query    string     true  "Name of city"
// @Success 200  {object} models.ShortForecastResponse "Short forecast"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Router /get_short_forecast [get]
func ShortPredGetter(log *slog.Logger, cities WeatherInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "/app_api/short_prediction_getter"
		city := r.URL.Query().Get("city")
		cities, err, stat := cities.GetShortPred(log, city)
		if err != nil {
			log.Error("Failed to get prediction", slog.Any("err", err), slog.String("path", path))
			w.WriteHeader(stat)
			return
		}
		log.Info("Get prediction success")
		render.JSON(w, r, cities)

	}
}
