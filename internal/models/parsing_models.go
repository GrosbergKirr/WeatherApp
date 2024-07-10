package models

type WeatherResponse struct {
	Cod     string   `json:"cod"`
	Message int      `json:"message"`
	Cnt     int      `json:"cnt"`
	List    []Record `json:"list"`
}

type Record struct {
	Dt         int64        `json:"dt"`
	Main       MainInfo     `json:"bin"`
	Weather    []WeatherGet `json:"weather"`
	Clouds     Clouds       `json:"clouds"`
	Wind       Wind         `json:"wind"`
	Visibility int          `json:"visibility"`
	Pop        float64      `json:"pop"`
	Rain       Rain         `json:"rain"`
	Sys        Sys          `json:"sys"`
	DtTxt      string       `json:"dt_txt"`
}

type MainInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type WeatherGet struct {
	ID          int    `json:"id"`
	Main        string `json:"bin"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Rain struct {
	ThreeH float64 `json:"3h"`
}

type Sys struct {
	Pod string `json:"pod"`
}
