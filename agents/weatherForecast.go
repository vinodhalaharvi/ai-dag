package agents

import (
	"ai-dag/config"
	"ai-dag/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type WeatherForecast struct {
}

func NewWeatherForecast() *WeatherForecast {
	return &WeatherForecast{}
}

func (owc *WeatherForecast) Do(
	config *config.DagConfig,
	agentId string,
	resultCh map[string]chan string,
	childResults map[string]string,
) {
	var weatherResponse *CurrentWeatherResponse
	format := "https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=%s&lang=%s&units=%s"
	parameters := config.Agents[agentId].QueryParameters
	appId := os.Getenv("OPEN_WEATHER_API_KEY")
	if appId == "" {
		fmt.Println("OPEN_WEATHER_API_KEY not set")
		return
	}
	url := fmt.Sprintf(
		format,
		parameters.Lat,
		parameters.Lon,
		appId,
		parameters.Lang,
		parameters.Units,
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return
	}
	jsonStr := utils.ToPrettyJsonFromObject(weatherResponse)
	resultCh[agentId] <- jsonStr
	close(resultCh[agentId])
}

type CurrentWeatherRequest struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Lang  string  `json:"lang"`
	Units string  `json:"units"`
}

type CurrentWeather struct {
	CurrentWeatherRequest
	CurrentWeatherResponse
}

type CurrentWeatherResponse struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Minutely []struct {
		Dt            int     `json:"dt"`
		Precipitation float64 `json:"precipitation"`
	} `json:"minutely"`
	Hourly []struct {
		Dt         int     `json:"dt"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Pop float64 `json:"pop"`
	} `json:"hourly"`
	Daily []struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Moonrise  int     `json:"moonrise"`
		Moonset   int     `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Summary   string  `json:"summary"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		WindGust  float64 `json:"wind_gust"`
		Weather   []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64 `json:"pop"`
		Uvi    float64 `json:"uvi"`
		Rain   float64 `json:"rain,omitempty"`
	} `json:"daily"`
	TemperatureAssessment string
}
