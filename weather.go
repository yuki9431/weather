package weather

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// 絶対零度 気温変換で使用する
const absoluteTmp = -273.15

type weatherInfo struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Rain  struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country  string `json:"country"`
		Timezone int    `json:"timezone"`
	} `json:"city"`
}

// TODO 都市を選べるようにする
func New() *weatherInfo {
	cityId := "1850147" // Tokyo
	appid := "63ef79e871474934c1bd707239475660"
	apiUrl := "http://api.openweathermap.org/data/2.5/forecast?id=" +
		cityId +
		"&" +
		"appid=" +
		appid

	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("天気情報の取得に失敗しました")
	}

	defer resp.Body.Close()

	// jsonデコード
	var weather weatherInfo
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		log.Fatal("jsonデコードに失敗しました")
	}
	return &weather
}

func (w *weatherInfo) GetCityName() string {
	return w.City.Name
}

func (w *weatherInfo) GetDates() []string {
	var dates []string
	for _, l := range w.List {
		dates = append(dates, l.DtTxt)
	}
	return dates
}

func (w *weatherInfo) GetDescriptions() []string {
	var descriptions []string
	for _, l := range w.List {
		for _, w := range l.Weather {
			descriptions = append(descriptions, w.Description)
		}
	}

	return descriptions
}

func (w *weatherInfo) GetTemps() []int {
	var maxTemps []int
	for _, l := range w.List {
		maxTemps = append(maxTemps, (int)(math.Round(l.Main.Temp+absoluteTmp)))
	}

	return maxTemps
}
