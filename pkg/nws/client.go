package nws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Peyton232/weather-forecast/pkg/model"
)

const (
	baseURL   = "https://api.weather.gov"
	userAgent = "weather-api-interview-example"
)

var ErrLocationUnsupported = errors.New("location not supported by NWS")

type Client interface {
	GetForecast(ctx context.Context, lat, lon float64) (*model.ForecastData, error)
}

type HTTPClient struct {
	httpClient *http.Client
}

func NewClient() Client {
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPClient) GetForecast(
	ctx context.Context,
	lat float64,
	lon float64,
) (*model.ForecastData, error) {

	pointURL := fmt.Sprintf("%s/points/%f,%f", baseURL, lat, lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pointURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrLocationUnsupported
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("points request failed: %s", resp.Status)
	}

	var pointResp model.PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointResp); err != nil {
		return nil, err
	}

	if pointResp.Properties.Forecast == "" {
		return nil, ErrLocationUnsupported
	}

	return c.fetchForecast(ctx, pointResp.Properties.Forecast)
}

func (c *HTTPClient) fetchForecast(
	ctx context.Context,
	url string,
) (*model.ForecastData, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast request failed: %s", resp.Status)
	}

	var forecastResp model.NWSForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, err
	}

	return &forecastResp.Properties, nil
}
