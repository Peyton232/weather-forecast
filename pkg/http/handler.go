package http

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/Peyton232/weather-forecast/pkg/nws"
	"github.com/Peyton232/weather-forecast/pkg/service"
)

type Handler struct {
	forecastService *service.ForecastService
	templates       *template.Template
}

func NewHandler(forecastService *service.ForecastService) *Handler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	return &Handler{
		forecastService: forecastService,
		templates:       tmpl,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func (h *Handler) GetForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqID := middleware.GetReqID(r.Context())

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	log.Printf("[reqID=%s] incoming forecast request lat=%s lon=%s", reqID, latStr, lonStr)

	if latStr == "" || lonStr == "" {
		log.Printf("[reqID=%s] missing lat or lon", reqID)
		http.Error(w, `{"error":"lat and lon are required"}`, http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Printf("[reqID=%s] invalid latitude: %v", reqID, err)
		http.Error(w, `{"error":"invalid latitude"}`, http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		log.Printf("[reqID=%s] invalid longitude: %v", reqID, err)
		http.Error(w, `{"error":"invalid longitude"}`, http.StatusBadRequest)
		return
	}

	// Basic coordinate validation
	if lat < -90 || lat > 90 {
		http.Error(w, `{"error":"latitude must be between -90 and 90"}`, http.StatusBadRequest)
		return
	}
	if lon < -180 || lon > 180 {
		http.Error(w, `{"error":"longitude must be between -180 and 180"}`, http.StatusBadRequest)
		return
	}

	forecast, err := h.forecastService.GetTodayForecast(r.Context(), lat, lon)
	if err != nil {

		if errors.Is(err, nws.ErrLocationUnsupported) {
			log.Printf("[reqID=%s] unsupported location lat=%f lon=%f", reqID, lat, lon)
			http.Error(w, `{"error":"forecast not available for this location"}`, http.StatusBadRequest)
			return
		}

		log.Printf("[reqID=%s] upstream forecast error: %v", reqID, err)
		http.Error(w, `{"error":"failed to retrieve forecast"}`, http.StatusBadGateway)
		return
	}

	log.Printf("[reqID=%s] forecast success lat=%f lon=%f", reqID, lat, lon)

	if err := json.NewEncoder(w).Encode(forecast); err != nil {
		log.Printf("[reqID=%s] response encoding error: %v", reqID, err)
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
