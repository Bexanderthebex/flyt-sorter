package main

import (
	"os"

	"github.com/Bexanderthebex/flight-sorter/api"
	"github.com/Bexanderthebex/flight-sorter/internal/sorter"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	zapLogger, _ := zap.NewProduction()

	e.Use(echozap.ZapLogger(zapLogger))
	e.Use(middleware.Recover())

	var activeSorter api.FlightSorter
	algo := os.Getenv("SORTING_ALGORITHM")
	switch algo {
	case "topological":
		activeSorter = sorter.TopologicalSort{}
	default:
		activeSorter = sorter.DegreeCountingSort{}
	}

	flightHandler := &api.FlightHandler{Sorter: activeSorter}

	e.POST("/calculate", flightHandler.SortFlights)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
