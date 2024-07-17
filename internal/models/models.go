package models

import (
	"time"
)

type City struct {
	Id      int     `json:"id" db:"id, omitempty"`
	Name    string  `json:"name" db:"name"`
	Country string  `json:"country" db:"country"`
	Lat     float64 `json:"lat" db:"latitude"`
	Lon     float64 `json:"lon" db:"longitude"`
}

type Forecast struct {
	Id           int       `json:"id" db:"id,omitempty"`
	CityName     string    `json:"city_name" db:"city_name"`
	Temp         float64   `json:"temp" db:"temperature"`
	Date         time.Time `json:"date" db:"date"`
	CityID       int       `json:"city_id" db:"city_id"`
	FullForecast []byte    `json:"full_forecast" db:"full_forecast,omitempty"`
}

type CitiesListResponse struct {
	Cities []string `json:"cities"`
}

type ShortForecastResponse struct {
	City     string      `json:"city"`
	Country  string      `json:"country"`
	MeanTemp float64     `json:"meanTemp"`
	Dates    []time.Time `json:"dates"`
}

type FullForecastRequest struct {
	City string `json:"city"`
	Date string `json:"date"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
