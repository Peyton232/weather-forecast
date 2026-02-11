package model

import "time"

// pointsResponse models /points/{lat},{lon}
type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// forecastResponse models /forecast
type NWSForecastResponse struct {
	Properties ForecastData `json:"properties"`
}

// ForecastData is the subset the service layer cares about.
type ForecastData struct {
	Periods []ForecastPeriod `json:"periods"`
}

type ForecastPeriod struct {
	Name          string    `json:"name"`
	StartTime     time.Time `json:"startTime"`
	Temperature   int       `json:"temperature"`
	ShortForecast string    `json:"shortForecast"`
}
