package base

import (
	"math/rand"
)

// Genetic Algorithm Parameters
var (
	mutationRate        float64 = 0.03
	elitism             bool    = true
	randomCrossoverRate         = false
	defCrossoverRate    float32 = 0.7
)

// var (
// 	mutationRate        float64 = 0.015
// 	elitism             bool    = true
// 	randomCrossoverRate         = true
// 	defCrossoverRate    float32 = 0.7
// )

func CrossoverRate() float32 {
	if randomCrossoverRate {
		return rand.Float32()
	}
	return defCrossoverRate
}

// Crossover : performs multi point cross over with 2 parents
// Assumption - parents have equal size
func Crossover(p1 Route, p2 Route) Route {
	// Size
	size := p1.RouteSize()
	// Child Trip
	c := Route{}
	c.InitRoute(size)

	// Number of crossover
	nc := int(CrossoverRate() * float32(size))
	if nc == 0 {
		// log.Println("no crossover")
		return p1
	}
	// Start positions of cross over for parent 1
	sp := int(rand.Float32() * float32(size))
	// End position of cross over for parent 1
	ep := (sp + nc) % size
	// Parent 2 slots
	p2s := make([]int, 0, size-nc)
	// log.Println(size, sp, nc, ep) // For debugging
	// Populate child with parent 1
	if sp < ep {
		for i := 0; i < size; i++ {
			if i >= sp && i < ep {
				c.SetPlace(i, p1.GetPlace(i))
			} else {
				p2s = append(p2s, i)
			}
		}
	} else if sp > ep {
		for i := 0; i < size; i++ {
			if !(i >= ep && i < sp) {
				c.SetPlace(i, p1.GetPlace(i))
			} else {
				p2s = append(p2s, i)
			}
		}
	}

	// For debugging
	// msPlace := ""
	j := 0
	// Populate child with parent 2 cities that are missing
	for i := 0; i < size; i++ {
		// Check if child contains Place
		if !c.ContainPlace(p2.GetPlace(i)) {
			c.SetPlace(p2s[j], p2.GetPlace(i))
			j++
			// For debugging
			// msPlace += p2.GetPlace(i).String() + " "
		}
	}
	return c
}

// Mutation : Performs swap mutation
// Chance of mutation for each Place based on mutation rate
func Mutation(in *Route) {
	// for each Place
	for p1 := 0; p1 < in.RouteSize(); p1++ {
		if rand.Float64() < mutationRate {
			// Select 2nd Place to perform swap
			p2 := int(float64(in.RouteSize()) * rand.Float64())
			// log.Println("Mutation occured", p1, "swap", p2)
			// Temp store Place
			c1 := in.GetPlace(p1)
			c2 := in.GetPlace(p2)
			// Swap Cities
			in.SetPlace(p1, c2)
			in.SetPlace(p2, c1)
		}
	}
}

// TripnamentSelection : select a group at random and pick the best parent
func TripnamentSelection(pop Population) Route {
	Routeny := Population{}
	// Routeny.InitEmpty(RouteSize)
	Routeny.InitEmpty(pop.GetRoute(0).RouteSize())

	for i := 0; i < pop.GetRoute(0).RouteSize(); i++ {
		r := int(rand.Float64() * float64(pop.PopulationSize()))
		Routeny.SaveRoute(i, *pop.GetRoute(r))
	}
	// fittest Trip
	fTrip := Routeny.GetFittest()
	return *fTrip
}

// EvolvePopulation : evolves population by :-
/*
	- Selecting 2 parents using Tripnament selection
	- Perform crossover to obtain child
	- Mutate child based on probability
	- return new population
*/
func EvolvePopulation(pop Population) Population {
	npop := Population{}
	npop.InitEmpty(pop.PopulationSize())

	popOffset := 0
	if elitism {
		npop.SaveRoute(0, *pop.GetFittest())
		popOffset = 1
	}

	for i := popOffset; i < npop.PopulationSize(); i++ {
		p1 := TripnamentSelection(pop)
		p2 := TripnamentSelection(pop)
		child := Crossover(p1, p2)
		Mutation(&child)
		npop.SaveRoute(i, child)
	}
	return npop
}
