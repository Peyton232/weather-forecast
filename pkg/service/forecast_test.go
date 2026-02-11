package service

import (
	"context"
	"testing"
	"time"

	"github.com/Peyton232/weather-forecast/pkg/model"
)

type mockNWSClient struct {
	data *model.ForecastData
	err  error
}

func (f *mockNWSClient) GetForecast(ctx context.Context, lat, lon float64) (*model.ForecastData, error) {
	return f.data, f.err
}

func TestGetTodayForecast_HappyPath(t *testing.T) {
	today := time.Now()

	fakeClient := &mockNWSClient{
		data: &model.ForecastData{
			Periods: []model.ForecastPeriod{
				{
					StartTime:     today,
					Temperature:   90,
					ShortForecast: "Sunny",
				},
			},
		},
	}

	service := NewForecastService(fakeClient)

	result, err := service.GetTodayForecast(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ShortForecast != "Sunny" {
		t.Errorf("expected Sunny, got %s", result.ShortForecast)
	}

	if result.TemperatureCategory != "hot" {
		t.Errorf("expected hot, got %s", result.TemperatureCategory)
	}
}

func TestCategorizeTemperature(t *testing.T) {
	tests := []struct {
		temp     int
		expected string
	}{
		{90, "hot"},
		{70, "moderate"},
		{40, "cold"},
	}

	for _, tt := range tests {
		result := categorizeTemperature(tt.temp)
		if result != tt.expected {
			t.Errorf("temp %d: expected %s, got %s", tt.temp, tt.expected, result)
		}
	}
}
