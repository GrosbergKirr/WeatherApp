package app_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

func GetWeather(log *slog.Logger,
	client http.Client,
	respStruct models.WeatherResponse,
	city models.City,
	weatherList *[]models.Weather,
	sideApiUrl string,
	apiKey string) (error, int) {
	urlBody := sideApiUrl + fmt.Sprintf("lat=%f&lon=%f&units=metric&appid=", city.Lat, city.Lon) + apiKey
	resp, err := client.Get(urlBody)
	if err != nil {
		log.Error("failed to get response", slog.String("error", err.Error()))
		return err, resp.StatusCode
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Read request body error")
		return err, resp.StatusCode
	}
	err = json.Unmarshal(respBody, &respStruct)
	if err != nil {
		log.Error("failed to decode json")
		return err, resp.StatusCode
	}
	defer resp.Body.Close()

	var weather models.Weather
	weather.FullForecast = respBody
	for prediction := range respStruct.List {
		dayTime := strings.Split(respStruct.List[prediction].DtTxt, " ")[1]
		if dayTime == "12:00:00" {
			weather.CityName = city.Name
			weather.Temp = respStruct.List[prediction].Main.Temp
			weather.Date = time.Unix(respStruct.List[prediction].Dt, 0)
			*weatherList = append(*weatherList, weather)
		}

	}

	return nil, resp.StatusCode

}
