package sorter

import (
	"sort"

	"github.com/Bexanderthebex/flight-sorter/pkg"
)

type Flight struct {
	Source      string
	Destination string
}

type FlightResult struct {
	Source       string
	Destinations []string
}

type TopologicalSort struct{}

func (ts TopologicalSort) Sort(flights []Flight) (FlightResult, error) {
	flightDependencyCount := make(map[string]int)
	flightGraph := make(map[string][]string)
	airports := make(map[string]bool)
	outDegree := make(map[string]int)

	for _, flight := range flights {
		flightGraph[flight.Source] = append(flightGraph[flight.Source], flight.Destination)
		flightDependencyCount[flight.Destination] += 1

		airports[flight.Source] = true
		airports[flight.Destination] = true
		outDegree[flight.Source]++
	}

	var flightEliminationQueue []string
	endNodesFound := 0

	// ensure all elements are in dependency count
	for airport := range airports {
		if _, ok := flightDependencyCount[airport]; !ok {
			flightDependencyCount[airport] = 0
		}
	}

	for airport := range airports {
		if flightDependencyCount[airport] == 0 {
			flightEliminationQueue = append(flightEliminationQueue, airport)
		}
		if outDegree[airport] == 0 {
			endNodesFound++
		}
	}

	flightPaths := len(flightEliminationQueue)
	if flightPaths <= 0 {
		return FlightResult{}, pkg.InvalidParameterError{Message: pkg.ErrMsgNoStartingPath, Code: pkg.ErrCodeNoStartingPath}
	}

	if flightPaths > 1 {
		return FlightResult{}, pkg.InvalidParameterError{Message: pkg.ErrMsgMultipleStartingPath, Code: pkg.ErrCodeMultipleStartingPath}
	}

	sortedFlights := make([]string, 0, len(flightDependencyCount))
	for len(flightEliminationQueue) > 0 {
		flight := flightEliminationQueue[0]
		flightEliminationQueue = flightEliminationQueue[1:]

		sortedFlights = append(sortedFlights, flight)

		for _, dependentFlight := range flightGraph[flight] {
			flightDependencyCount[dependentFlight] -= 1
			if flightDependencyCount[dependentFlight] == 0 {
				flightEliminationQueue = append(flightEliminationQueue, dependentFlight)
			}
		}

		delete(flightDependencyCount, flight)
	}

	var destinations []string
	for airport := range airports {
		if outDegree[airport] == 0 {
			destinations = append(destinations, airport)
		}
	}
	sort.Strings(destinations)

	return FlightResult{Source: sortedFlights[0], Destinations: destinations}, nil
}

type DegreeCountingSort struct{}

func (dcs DegreeCountingSort) Sort(flights []Flight) (FlightResult, error) {
	if len(flights) == 0 {
		return FlightResult{}, pkg.InvalidParameterError{Message: pkg.ErrMsgNoStartingPath, Code: pkg.ErrCodeNoStartingPath}
	}

	inDegree := make(map[string]int)
	outDegree := make(map[string]int)
	airports := make(map[string]bool)

	for _, flight := range flights {
		outDegree[flight.Source]++
		inDegree[flight.Destination]++
		airports[flight.Source] = true
		airports[flight.Destination] = true
	}

	var startNode string
	var endNodes []string
	startNodesFound := 0

	for airport := range airports {
		in := inDegree[airport]
		out := outDegree[airport]

		if in == 0 {
			startNode = airport
			startNodesFound++
		}
		if out == 0 {
			endNodes = append(endNodes, airport)
		}
	}

	if startNodesFound == 0 {
		return FlightResult{}, pkg.InvalidParameterError{Message: pkg.ErrMsgNoStartingPath, Code: pkg.ErrCodeNoStartingPath}
	}
	if startNodesFound > 1 {
		return FlightResult{}, pkg.InvalidParameterError{Message: pkg.ErrMsgMultipleStartingPath, Code: pkg.ErrCodeMultipleStartingPath}
	}

	sort.Strings(endNodes)
	return FlightResult{Source: startNode, Destinations: endNodes}, nil
}
