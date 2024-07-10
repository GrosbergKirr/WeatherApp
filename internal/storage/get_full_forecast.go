package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func (s *Storage) GetFullPred(log *slog.Logger, param models.Weather) (models.Record, error, int) {
	const path = "/storage/get_full_forecast"
	var fullPredJson []byte
	queryCity := "SELECT full_forecast FROM weather where city_name = $1 and date = $2"
	if err := s.Db.Get(&fullPredJson, queryCity, param.CityName, param.Date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to get cities from db. No forecast for this parameters. (Valid date layout: '2006-01-02 15:04:05')", slog.String("path", path))
			return models.Record{}, errors.New("no forecast for this city"), http.StatusBadRequest
		}
		log.Error("Failed to get cities from db", slog.String("path", path))
		return models.Record{}, err, http.StatusInternalServerError
	}
	var fullPredStruct models.WeatherResponse
	if err := json.Unmarshal(fullPredJson, &fullPredStruct); err != nil {
		log.Error("failed to decode data from db", slog.String("path", path))
		return models.Record{}, err, http.StatusInternalServerError
	}

	layout := "2006-01-02 15:04:05"
	for i := range fullPredStruct.List {
		dateTime, err := time.Parse(layout, fullPredStruct.List[i].DtTxt)
		if err != nil {
			log.Error("Can't convert request string to DateTime", slog.String("path", path))
			return models.Record{}, err, http.StatusBadRequest
		}
		if dateTime == param.Date {
			return fullPredStruct.List[i], nil, http.StatusOK
		}
	}
	return models.Record{}, nil, http.StatusInternalServerError
}
