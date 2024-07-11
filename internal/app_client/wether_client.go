package app_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func GetWeather(log *slog.Logger,
	client http.Client,
	mu *sync.Mutex,
	city models.City,
	weatherList *[]models.Forecast,
	sideApiUrl string,
	apiKey string) (int, error) {
	u, err := url.Parse(sideApiUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return http.StatusInternalServerError, err
	}
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", city.Lat))
	q.Set("lon", fmt.Sprintf("%f", city.Lat))
	q.Set("units", "metric")
	q.Set("appid", apiKey)
	u.RawQuery = q.Encode()
	urlBody := u.String()

	resp, err := client.Get(urlBody)
	if err != nil {
		log.Error("failed to get response", slog.String("error", err.Error()))
		return resp.StatusCode, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Read request body error")
		return resp.StatusCode, err
	}
	var respStruct models.WeatherResponse
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		log.Error("failed to decode json")
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	var weather models.Forecast
	weather.FullForecast = respBody
	for _, prediction := range respStruct.List {
		dayTime := strings.Split(prediction.DtTxt, " ")[1]
		if dayTime == "12:00:00" {
			weather.CityName = city.Name
			weather.Temp = prediction.Main.Temp
			weather.Date = time.Unix(prediction.Dt, 0)
			mu.Lock()
			*weatherList = append(*weatherList, weather)
			mu.Unlock()
		}

	}
	return resp.StatusCode, err
}
