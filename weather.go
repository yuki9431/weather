package weather

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"
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
func New(cityId string, appid string) *weatherInfos {
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

func (w *weatherInfos) GetIcons() []string {
	var icons []string
	for _, l := range w.List {
		for _, lw := range l.Weather {
			icons = append(icons, lw.Icon)
		}
	}

	return icons
}

func (w *weatherInfos) GetDates() []time.Time {
	var times []time.Time

	for _, l := range w.List {
		date, _ := time.Parse("2006-01-02 15:04:05", l.DtTxt)
		times = append(times, date)
	}

	return times
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

func (w *weatherInfos) ConvertIconToWord(icon string) string {
	var word string

	switch icon {
	case "01d", "01n":
		word = "快晴"
	case "02d", "02n":
		word = "晴れ"
	case "03d", "04d", "03n", "04n":
		word = "曇り"
	case "09d", "09n":
		word = "小雨"
	case "10d", "10n":
		word = "雨"
	case "11d", "11n":
		word = "雷雨"
	case "13d", "13n":
		word = "雪"
	case "50d", "50n":
		word = "霧"
	default:
		word = "該当情報無し"

	}
	return word
}

func (w *weatherInfos) GetInfoFromDate(target time.Time) *weatherInfos {
	const (
		layoutWeatherDate = "2006-01-02 15:04:05" // => YYYY-MM-DD hh:dd:ss
		layout            = "2006-01-02"          // => YYYY-MM-DD
	)
	var weatherInfosToday weatherInfos
	weatherInfosToday.City.Name = w.City.Name

	for i, date := range w.GetDates() {
		if t := date; target.Format(layout) == t.Format(layout) {
			weatherInfosToday.List = append(weatherInfosToday.List, w.List[i])
		}
	}
	return &weatherInfosToday
}
