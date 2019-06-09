package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

// 絶対零度 気温変換で使用する
const absoluteTmp = -273.15

type weatherInfos struct {
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
}

// TODO 都市を選べるようにする
func New() *weatherInfos {
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
	var weather weatherInfos
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		log.Fatal("jsonデコードに失敗しました")
	}
	return &weather
}

func (w *weatherInfos) GetCityName() string {
	return w.City.Name
}

func (w *weatherInfos) GetIcon() []string {
	var icon []string
	for _, l := range w.List {
		for _, w := range l.Weather {
			icon = append(icon, w.Icon)
		}
	}

	return icon
}

func (w *weatherInfos) GetDates() []string {
	var dates []string
	for _, l := range w.List {
		dates = append(dates, l.DtTxt)
	}
	return dates
}

func (w *weatherInfos) GetDescriptions() []string {
	var descriptions []string
	for _, l := range w.List {
		for _, w := range l.Weather {
			descriptions = append(descriptions, w.Description)
		}
	}

	return descriptions
}

func (w *weatherInfos) GetTemps() []int {
	var maxTemps []int
	for _, l := range w.List {
		maxTemps = append(maxTemps, (int)(math.Round(l.Main.Temp+absoluteTmp)))
	}

	return maxTemps
}

func main() {
	weather := New()
	fmt.Println(weather.GetCityName())
	for i, date := range weather.GetDates() {
		fmt.Println("date: " + (string)(date))
		fmt.Println("天気: " + (string)(weather.GetDescriptions()[i]))
		fmt.Println(weather.GetTemps()[i])

	}

}
