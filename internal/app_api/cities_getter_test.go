package app_api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	_ "strconv"
	"testing"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	cfg    = internal.SetupConfig("/home/kirrgross/GolandProjects/Weather_app/config/config.yaml")
	log, _ = internal.SetupLogger(cfg)
	mockDB = &MockDatabase{}
)

// Мокаем бдшку + интерфейс + методы
type MockDatabase struct{}

func (m *MockDatabase) GetCities(log *slog.Logger, perPage, offset int) (models.CitiesListResponse, int, error) {
	return models.CitiesListResponse{
		//реальность
		Cities: []string{
			"New York",
			"London",
			"Paris",
		},
	}, http.StatusOK, nil
}

func (m *MockDatabase) GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, int, error) {
	return models.ShortForecastResponse{}, 0, nil
}

func (m *MockDatabase) GetFullPred(log *slog.Logger, param models.Forecast) (models.Record, int, error) {
	return models.Record{}, 0, nil
}

func TestCitiesGetter_Success(t *testing.T) {
	// ожидание
	expected := `{"cities": ["New York","London","Paris"]}`
	// разворачиваем обертку над стандартным хэндлером
	handler := CitiesGetter(log, mockDB)

	req := httptest.NewRequest(http.MethodGet, "/app_api/cities_getter?page=1&per_page=10", nil)
	w := httptest.NewRecorder()

	// дергаем хэндлер
	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	// чекаем ответ и код
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

type MockDatabaseWithError struct{}

func (m *MockDatabaseWithError) GetCities(log *slog.Logger, perPage, offset int) (models.CitiesListResponse, int, error) {
	return models.CitiesListResponse{}, http.StatusInternalServerError, errors.New("database error")
}

func (m *MockDatabaseWithError) GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, int, error) {
	return models.ShortForecastResponse{}, 0, nil
}

func (m *MockDatabaseWithError) GetFullPred(log *slog.Logger, param models.Forecast) (models.Record, int, error) {
	return models.Record{}, 0, nil
}

func TestCitiesGetter_Error(t *testing.T) {
	mockDBWithError := &MockDatabaseWithError{}
	handler := CitiesGetter(log, mockDBWithError)

	req := httptest.NewRequest(http.MethodGet, "/app_api/cities_getter?page=1&per_page=10", nil)
	w := httptest.NewRecorder()

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}
