package storage

import (
	"errors"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestGetShortPred_Success(t *testing.T) {
	var (
		city = "city"
	)
	expectedCountry := "country"
	expectedWeather := models.Forecast{
		Id:       1,
		CityName: "name",
		Temp:     111.1,
		Date:     time.Now(),
		CityID:   1,
	}
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	cfg := internal.SetupConfig("/home/kirrgross/GolandProjects/Weather_app/config/config.yaml")
	log, err := internal.SetupLogger(cfg)

	const query = "SELECT country FROM cities where name = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(city).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"country"}).
				AddRow(expectedCountry)).
		RowsWillBeClosed()

	const queryWeather = "SELECT temperature, date FROM weather where city_name = $1"
	mock.ExpectQuery(regexp.QuoteMeta(queryWeather)).
		WithArgs(city).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"temperature", "date"}).
				AddRow(expectedWeather.Temp, expectedWeather.Date)).
		RowsWillBeClosed()

	storage := Storage{Db: sqlxDB}
	rows, err, status := storage.GetShortPred(log, city)
	require.NoError(t, err)
	require.NotNil(t, rows)
	require.Equal(t, status, http.StatusOK)
	require.Equal(t, expectedCountry, rows.Country)

}

func TestGetShortPred_Error(t *testing.T) {
	var (
		city = "city"
	)
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	cfg := internal.SetupConfig("/home/kirrgross/GolandProjects/Weather_app/config/config.yaml")
	log, err := internal.SetupLogger(cfg)

	const query = "SELECT country FROM cities where name = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(city).
		WillReturnError(errors.New("error"))

	const queryWeather = "SELECT temperature, date FROM weather where city_name = $1"
	mock.ExpectQuery(regexp.QuoteMeta(queryWeather)).
		WithArgs(city).
		WillReturnError(errors.New("error"))

	storage := Storage{Db: sqlxDB}
	rows, err, status := storage.GetShortPred(log, city)
	require.Error(t, err)
	require.Zero(t, rows.MeanTemp)
	require.Nil(t, rows.Dates)
	require.Equal(t, status, http.StatusInternalServerError)
}
