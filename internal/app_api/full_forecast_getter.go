package app_api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/go-chi/render"
)

// FullPredGetter godoc
// @Summary Get list of available cities
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   City&Date  body		models.FullForecastRequest    true  "Name of city and date of full forecast" example({"city": "Minsk", "date": "2024-07-12 12:00:00"})
// @Success 200  {object} models.Record "Full forecast"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
// @Router /get_full_forecast [post]
func FullPredGetter(log *slog.Logger, record WeatherInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path string = "app_api/full_forecast_getter"
		var req models.FullForecastRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail to decode json. Valid date layout: \"2006-01-02 15:04:05\"", slog.Any("err: ", err), slog.String("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(req.Date)
		log.Info("Get and decode json success")

		var forecast models.Forecast
		forecast.CityName = req.City
		layout := "2006-01-02 15:04:05"
		date, _ := time.Parse(layout, req.Date)
		fmt.Println(date)
		forecast.Date = date
		record, err, stat := record.GetFullPred(log, forecast)
		if err != nil {
			log.Error("Failed to get full prediction", slog.Any("err", err), slog.String("path", path))
			w.WriteHeader(stat)
			return
		}
		render.JSON(w, r, record)
	}
}
