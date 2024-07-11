package storage

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestGetCities_Success(t *testing.T) {
	var (
		offset  = 10
		perPage = 5
	)

	expectedCity := models.City{
		Id:      1,
		Name:    "name",
		Country: "country",
		Lat:     2,
		Lon:     3,
	}

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	log := internal.SetupLogger()

	const query = `SELECT name FROM cities LIMIT $1 OFFSET $2`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(perPage, offset).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"name"}).
				AddRow(expectedCity.Name)).
		RowsWillBeClosed()

	storage := Storage{Db: sqlxDB}

	rows, err, status := storage.GetCities(log, perPage, offset)

	require.NoError(t, err)
	require.NotNil(t, rows)
	require.Equal(t, status, http.StatusOK)
	require.Equal(t, expectedCity.Name, rows.Cities[0])
}

func TestGetCities_Error(t *testing.T) {
	var (
		offset  = 10
		perPage = 5
	)

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	log := internal.SetupLogger()

	const query = `SELECT name FROM cities LIMIT $1 OFFSET $2`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(perPage, offset).
		WillReturnError(errors.New("error"))

	storage := Storage{Db: sqlxDB}

	rows, err, status := storage.GetCities(log, perPage, offset)

	require.Error(t, err)
	require.Nil(t, rows.Cities)
	require.Equal(t, status, http.StatusInternalServerError)
}
