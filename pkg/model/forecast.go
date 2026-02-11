package model

// ForecastResponse is returned to API consumers.
type ForecastResponse struct {
	ShortForecast       string `json:"shortForecast"`
	Temperature         int    `json:"temperature"`
	TemperatureCategory string `json:"temperatureCategory"`
}
