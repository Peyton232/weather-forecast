package service

import (
	"context"
	"errors"
	"time"

	"github.com/Peyton232/weather-forecast/pkg/model"
	"github.com/Peyton232/weather-forecast/pkg/nws"
)

// ForecastService contains business logic only.
type ForecastService struct {
	nwsClient nws.Client
}

// NewForecastService wires the dependency.
func NewForecastService(nwsClient nws.Client) *ForecastService {
	return &ForecastService{
		nwsClient: nwsClient,
	}
}

// GetTodayForecast  retrieves today's forecast and applies business rules.
func (s *ForecastService) GetTodayForecast(
	ctx context.Context,
	lat float64,
	lon float64,
) (*model.ForecastResponse, error) {

	forecast, err := s.nwsClient.GetForecast(ctx, lat, lon)
	if err != nil {
		return nil, err
	}

	period, err := selectTodayPeriod(forecast.Periods)
	if err != nil {
		return nil, err
	}

	return &model.ForecastResponse{
		ShortForecast:       period.ShortForecast,
		Temperature:         period.Temperature,
		TemperatureCategory: categorizeTemperature(period.Temperature),
	}, nil
}

// --- helpers ---
func selectTodayPeriod(periods []model.ForecastPeriod) (*model.ForecastPeriod, error) {
	if len(periods) == 0 {
		return nil, errors.New("no forecast periods returned")
	}

	today := time.Now().Format("2006-01-02")

	for _, p := range periods {
		if p.StartTime.Format("2006-01-02") == today {
			return &p, nil
		}
	}

	// Fallback: assume first period is today
	return &periods[0], nil
}

func categorizeTemperature(temp int) string {
	switch {
	case temp >= 85:
		return "hot"
	case temp <= 50:
		return "cold"
	default:
		return "moderate"
	}
}
