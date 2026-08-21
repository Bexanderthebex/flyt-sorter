package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bexanderthebex/flight-sorter/internal/sorter"
	"github.com/Bexanderthebex/flight-sorter/pkg"
	"github.com/labstack/echo/v4"
)

type FlightSorter interface {
	Sort(flights []sorter.Flight) (sorter.FlightResult, error)
}

type FlightHandler struct {
	Sorter FlightSorter
}

func (h *FlightHandler) SortFlights(ctx echo.Context) error {
	var flightInput flights
	if err := ctx.Bind(&flightInput); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input format. Should be an array of 2 element strings",
			"code":    http.StatusBadRequest,
		})
	}

	flightsToSort := make([]sorter.Flight, 0, len(flightInput))
	for _, f := range flightInput {
		inputLength := len(f)
		if inputLength != 2 {
			return ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": fmt.Sprintf("invalid input format, expecting 2 array elements but got %d", inputLength),
				"code":    http.StatusBadRequest,
			})
		}

		flightsToSort = append(flightsToSort, sorter.Flight{Source: strings.ToUpper(f[0]), Destination: strings.ToUpper(f[1])})
	}

	sortedFlightPath, err := h.Sorter.Sort(flightsToSort)
	if err != nil {
		var invalidParamErr pkg.InvalidParameterError
		if errors.As(err, &invalidParamErr) {
			return ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": invalidParamErr.Message,
				"code":    invalidParamErr.Code,
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, []any{sortedFlightPath.Source, sortedFlightPath.Destinations})
}
