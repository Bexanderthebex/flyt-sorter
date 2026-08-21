package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Bexanderthebex/flight-sorter/api"
	"github.com/Bexanderthebex/flight-sorter/internal/sorter"
	"github.com/Bexanderthebex/flight-sorter/pkg"
	"github.com/labstack/echo/v4"
)

func TestIntegration_CalculateEndpoint_InMemory(t *testing.T) {
	e := echo.New()

	topologicalSorter := &sorter.TopologicalSort{}
	flightHandler := &api.FlightHandler{Sorter: topologicalSorter}

	e.POST("/calculate", flightHandler.SortFlights)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResp   any
	}{
		{
			name:           "HappyPath",
			requestBody:    `[["SFO", "EWR"], ["ATL", "EWR"], ["SFO", "ATL"]]`,
			expectedStatus: http.StatusOK,
			expectedResp:   []any{"SFO", []any{"EWR"}},
		},
		{
			name:           "BadCycleRequestDetected",
			requestBody:    `[["A", "B"], ["B", "A"]]`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": pkg.ErrMsgNoStartingPath,
				"code":    string(pkg.ErrCodeNoStartingPath),
			},
		},
		{
			name:           "InvalidInputFormatJSON",
			requestBody:    `{ "SF": "VE" }`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "invalid input format. Should be an array of 2 element strings",
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name:           "InvalidInputFormatWrongLength",
			requestBody:    `[["AB", "C", "D"]]`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "invalid input format, expecting 2 array elements but got 3",
				"code":    float64(http.StatusBadRequest),
			},
		},
		{
			name:           "NodesAreCaseInsensitive",
			requestBody:    `[["sfo", "EWR"], ["atL", "EWR"], ["sfO", "Atl"]]`,
			expectedStatus: http.StatusOK,
			expectedResp:   []any{"SFO", []any{"EWR"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, rec.Code)
			}

			var response any
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to decode response: %v. Body was: %s", err, rec.Body.String())
			}

			if !reflect.DeepEqual(response, tt.expectedResp) {
				t.Errorf("Expected response %v, got %v", tt.expectedResp, response)
			}
		})
	}
}
