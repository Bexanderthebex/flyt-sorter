# flight-sorter
__Story__: There are over 100,000 flights a day, with millions of people and cargo being transferred around the world. With so many people and different carrier/agency groups, it can be hard to track where a person might be. In order to determine the flight path of a person, we must sort through all of their flight records.

__Goal__: To create a simple microservice API that can help us understand and track how a particular person's flight path may be queried. The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.

Required JSON structure: 
- [["SFO", "EWR"]]                                                 => ["SFO", "EWR"]
- [["ATL", "EWR"], ["SFO", "ATL"]]                                 => ["SFO", "EWR"]
- [["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]

### Assumptions

In order for the solution to be useful, it needs to be deterministic. For example, when presented with this input, will the program be able to clearly tell what is the exact starting path?

```json
[
    ["SFO", "ATL"],
    ["ATL", "JFK"],
    ["JFK", "SFO"]
]
```

In order for the program to be useful, it will need to eliminate ambiguity and only deal with problem that it knows that a clear answer can be given. The following inputs will be discarded:

1. No possible starting paths

2. Multiple possible starting paths

The program also explicitly supports **branching flight paths that result in multiple final destinations**. In real-world logistics (like tracking cargo that splits at a layover hub), having multiple endpoints is a highly critical edge case to solve. To support this, the API return format has been slightly adapted from the original specification to return an array of destinations: `["starting_node", ["ending_node_1", "ending_node_2"]]`.

Aside from that, the strings will be treated in a case insesitive manner as they represent the same semantic meaning in the real world. E.g. `SFO` is the same as `sfo`.

### Solution

It is quite clear that this is a sorting problem so the initial thought is to reach out for Kahn's algorithm. However, since the problem only requires finding the absolute starting and ending paths without needing to reconstruct the exact layover sequence, Kahn's algorithm might not be optimal since it requires building and traversing a heavy dependency graph.

Instead, this problem can be elegantly solved using **Degree Counting** (a concept rooted in Eulerian Paths). By simply tallying the "inbound" (arrivals) and "outbound" (departures) flights for each airport, we can mathematically guarantee the endpoints without any path traversal:
* The **Starting Airport** will be the *only* node with 0 inbound flights (`in == 0`).
* The **Final Destinations** will be the *only* nodes with 0 outbound flights (`out == 0`).

This completely bypasses the need to build a complex adjacency list or maintain a processing queue. While both Kahn's algorithm and Degree Counting are technically `O(N)` in time complexity (or `O(V + E)`), Degree Counting boasts a much smaller constant factor and is significantly more memory-efficient (as proven by our benchmarks) because it avoids the heavy slice allocations required to build a graph.

### How to run?

#### Prerequisites:
The program assumes that there is docker installed

First, run `go mod download`

To run the tests:

`make test`

To run the benchmark:

`make bench-test`

To run in docker:

`make docker-run`

**Note** The program is meant to run in port 8080 of the host machine. If in any case the program errors out, it might be because there's other programs might be using that port. In that case, specify a different port by

`make docker-run PORT=<desired_port>`

### API Reference

**Endpoint:** `POST /calculate`

**Headers:**
- `Content-Type: application/json`

**Request Body:**
A JSON array of flight segments, where each segment is a 2-element array of strings `["Source", "Destination"]`.
```json
[
  ["ATL", "EWR"],
  ["SFO", "ATL"]
]
```

**Success Response (200 OK):**
A JSON array containing the calculated starting airport and an array of all possible final destinations:
```json
[
  "SFO",
  [
    "EWR"
  ]
]
```
