package sorter

import (
	"fmt"
	"testing"

	"github.com/Bexanderthebex/flight-sorter/pkg"
)

func TestSorters(t *testing.T) {
	testCases := []struct {
		name        string
		input       []Flight
		expected    FlightResult
		expectedErr error
	}{
		{
			name: "HappyPath",
			input: []Flight{
				{Source: "IND", Destination: "EWR"},
				{Source: "SFO", Destination: "ATL"},
				{Source: "GSO", Destination: "IND"},
				{Source: "ATL", Destination: "GSO"},
			},
			expected: FlightResult{Source: "SFO", Destinations: []string{"EWR"}},
		},
		{
			name: "DivergeButConverge",
			input: []Flight{
				{Source: "SFO", Destination: "ATL"},
				{Source: "ATL", Destination: "GSO"},
				{Source: "ATL", Destination: "IND"},
				{Source: "GSO", Destination: "WI"},
				{Source: "IND", Destination: "WI"},
				{Source: "WI", Destination: "NY"},
			},
			expected: FlightResult{Source: "SFO", Destinations: []string{"NY"}},
		},
		{
			name: "DivergeButConvergeWithTwoOpenNodes",
			input: []Flight{
				{Source: "SFO", Destination: "ATL"},
				{Source: "ATL", Destination: "GSO"},
				{Source: "ATL", Destination: "IND"},
				{Source: "GSO", Destination: "WI"},
				{Source: "WI", Destination: "NY"},
			},
			expected: FlightResult{Source: "SFO", Destinations: []string{"IND", "NY"}},
		},
		{
			name: "NoPossibleStartingPath",
			input: []Flight{
				{Source: "SFO", Destination: "ATL"},
				{Source: "ATL", Destination: "JFK"},
				{Source: "JFK", Destination: "SFO"},
			},
			expectedErr: pkg.InvalidParameterError{Message: pkg.ErrMsgNoStartingPath, Code: pkg.ErrCodeNoStartingPath},
		},
		{
			name: "MultipleStartingPaths",
			input: []Flight{
				{Source: "SFO", Destination: "ATL"},
				{Source: "IND", Destination: "GSO"},
			},
			expectedErr: pkg.InvalidParameterError{Message: pkg.ErrMsgMultipleStartingPath, Code: pkg.ErrCodeMultipleStartingPath},
		},
		{
			name: "MultipleEndingDestinations",
			input: []Flight{
				{Source: "SFO", Destination: "ATL"},
				{Source: "ATL", Destination: "GSO"},
				{Source: "ATL", Destination: "IND"},
				{Source: "GSO", Destination: "NY"},
			},
			expected: FlightResult{Source: "SFO", Destinations: []string{"IND", "NY"}},
		},
	}

	type Sorter interface {
		Sort(flights []Flight) (FlightResult, error)
	}

	sorters := []struct {
		name   string
		sorter Sorter
	}{
		{"TopologicalSort", &TopologicalSort{}},
		{"DegreeCountingSort", &DegreeCountingSort{}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			for _, s := range sorters {
				t.Run(s.name, func(t *testing.T) {
					flight, err := s.sorter.Sort(test.input)

					if test.expectedErr != nil {
						if test.expectedErr.Error() != err.Error() {
							t.Error(fmt.Sprintf("expected %q got %q", test.expectedErr.Error(), err.Error()))
						}
						return
					}

					if test.expectedErr == nil && err != nil {
						t.Error(fmt.Sprintf("unexpected error: %s", err.Error()))
					}

					if test.expected.Source != flight.Source {
						t.Error(fmt.Sprintf("expected source %s got %s", test.expected.Source, flight.Source))
					}
					
					if len(test.expected.Destinations) != len(flight.Destinations) {
						t.Errorf("expected %d destinations got %d", len(test.expected.Destinations), len(flight.Destinations))
					} else {
						for i := range test.expected.Destinations {
							if test.expected.Destinations[i] != flight.Destinations[i] {
								t.Errorf("expected destination %s got %s", test.expected.Destinations[i], flight.Destinations[i])
							}
						}
					}
				})
			}
		})
	}
}
