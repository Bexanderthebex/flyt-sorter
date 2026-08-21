package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bexanderthebex/flight-sorter/api"
	"github.com/Bexanderthebex/flight-sorter/internal/sorter"
	"github.com/labstack/echo/v4"
)

// generateFlights creates a long linear flight path of N flights.
// e.g., ["Airport0", "Airport1"], ["Airport1", "Airport2"]...
func generateFlights(n int) string {
	flights := make([][2]string, 0, n)
	for i := 0; i < n; i++ {
		flights = append(flights, [2]string{fmt.Sprintf("Airport%d", i), fmt.Sprintf("Airport%d", i+1)})
	}

	b, _ := json.Marshal(flights)
	return string(b)
}

func setupEchoWithAlgorithm(algo string) *echo.Echo {
	e := echo.New()

	var activeSorter api.FlightSorter
	switch algo {
	case "degree":
		activeSorter = sorter.DegreeCountingSort{}
	default:
		activeSorter = sorter.TopologicalSort{}
	}

	flightHandler := &api.FlightHandler{Sorter: activeSorter}
	e.POST("/calculate", flightHandler.SortFlights)

	return e
}

func benchmarkSorting(b *testing.B, algo string, flightCount int) {
	e := setupEchoWithAlgorithm(algo)

	requestBody := generateFlights(flightCount)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}
}

func BenchmarkTopologicalSort_100(b *testing.B) {
	benchmarkSorting(b, "topological", 100)
}

func BenchmarkDegreeCountingSort_100(b *testing.B) {
	benchmarkSorting(b, "degree", 100)
}

func BenchmarkTopologicalSort_10000(b *testing.B) {
	benchmarkSorting(b, "topological", 10000)
}

func BenchmarkDegreeCountingSort_10000(b *testing.B) {
	benchmarkSorting(b, "degree", 10000)
}

func BenchmarkTopologicalSort_100000(b *testing.B) {
	benchmarkSorting(b, "topological", 100000)
}

func BenchmarkDegreeCountingSort_100000(b *testing.B) {
	benchmarkSorting(b, "degree", 100000)
}
