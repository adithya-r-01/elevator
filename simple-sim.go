package main

import (
	"flag"
	"fmt"
)

type Building struct {
	numFloors  int
	totalUsers int
	floors     []float64
}

func main() {
	var (
		numCycles    = flag.Int("num_simulations", 100, "Number of simulations to run")
		numFloors    = flag.Int("num_floors", 10, "Number of floors above ground (excluding ground)")
		initialUsers = flag.Int("initial_users", 111, "Initial number of users in the building")
		trickleUsers = flag.Int("trickle_users", 10, "Number of users to trickle in each cycle")
	)
	flag.Parse()

	n := *numFloors
	building := Building{
		floors: make([]float64, n+1),
	}

	initialDist := float64(*initialUsers) / float64(n+1)
	for i := 0; i <= n; i++ {
		building.floors[i] = initialDist
	}

	P := buildTransitionMatrix(n)

	for i := 0; i < *numCycles; i++ {
		building.floors = transition(building.floors, P, float64(*trickleUsers))
	}

	for i, users := range building.floors {
		if i == 0 {
			fmt.Printf("Ground Floor: %.2f users\n", users)
		} else {
			fmt.Printf("Floor %d: %.2f users\n", i, users)
		}
	}
}

func buildTransitionMatrix(n int) [][]float64 {
	size := n + 1
	P := make([][]float64, size)
	for i := range P {
		P[i] = make([]float64, size)
	}

	for i := 1; i <= n; i++ {
		P[i][0] = 0.9

		for j := 1; j <= n; j++ {
			if j != i {
				P[i][j] = 0.1 / float64(n-1)
			}
		}
	}

	for j := 1; j <= n; j++ {
		P[0][j] = 1.0 / float64(n)
	}

	return P
}

func transition(floors []float64, P [][]float64, trickle float64) []float64 {
	n := len(floors)
	newDist := make([]float64, n)

	for i := range n {
		for j := range n {
			newDist[j] += floors[i] * P[i][j]
		}
	}

	newDist[0] += trickle

	return newDist
}
