package base

import (
	"math"
	"math/rand"
	"sort"
)

// Route : route struct for optimisation
type Route struct {
	places   []Place
	fitness  float64
	distance float64
	walkDist float64
	time     float64
	walkTime float64
	match    float64
}

// InitRoute : Initialize Route with cities arranged randomly
func (r *Route) InitRoute(numberOfPlaces int) {
	r.places = make([]Place, numberOfPlaces)
}

// InitRoutePlaces init route
func (r *Route) InitRoutePlaces(tm TripManager) {
	r.InitRoute(tm.NumberOfPlaces())
	// Add all destination cities from RouteManager to Route
	for i := 0; i < tm.NumberOfPlaces(); i++ {
		r.SetPlace(i, tm.GetPlace(i))
	}
	r.places = ShufflePlaces(r.places)
}

func (a *Route) SetTime(places []Place) {
	// run tsp (no graph max)
	// r := Route{}

	// t.totalTime
	return
}

func (a *Route) SetMatch(cat []float64, mood []float64, pop float64) {
	sum := 0.0
	for i := 0; i < len(a.places); i++ {
		a.places[i].MC = a.places[i].GetPlaceMC(cat, mood, pop)
		sum += a.places[i].MC
	}
	a.match = sum
	// a.match = sum / float64(len(a.places))
}

func (a *Route) SetMatchCat(cat []float64) {
	sum := 0.0
	for i := 0; i < len(a.places); i++ {
		a.places[i].MC = a.places[i].getPlaceMCCat(cat)
		sum += a.places[i].MC
	}
	a.match = sum
}

func (a *Route) SetMatchMood(mood []float64) {
	sum := 0.0
	for i := 0; i < len(a.places); i++ {
		a.places[i].MC = a.places[i].getPlaceMCMood(mood)
		sum += a.places[i].MC
	}
	a.match = sum
}

func (a *Route) SetMatchPop(pop float64) {
	sum := 0.0
	for i := 0; i < len(a.places); i++ {
		a.places[i].MC = a.places[i].getPlaceMCPop(pop)
		sum += a.places[i].MC
	}
	a.match = sum
}

// GetMatch : Get matching coefficient of the combination of sights
func (r *Route) GetMatch() float64 {
	return r.match
}

// GetPlace : Get Place based on position in slice
func (r *Route) GetPlace(idx int) Place {
	return r.places[idx]
}

// GetPlaces : Get all Places
func (r *Route) GetPlaces() []Place {
	return r.places
}

// GetCats : Get Cats for route
func (r *Route) GetCats() []float64 { // how to improve aggregation of place cats / moods / pop into trip?
	cats := make([]float64, len(r.places[0].Cat))
	for i := 0; i < len(cats); i++ {
		for j := 0; j < len(r.places); j++ {
			cats[i] += r.places[j].Cat[i]
		}
		cats[i] = cats[i] / math.Sqrt(float64(len(r.places)))
	}
	return cats
}

// GetMoods : Get Moods for route
func (r *Route) GetMoods() []float64 { // how to improve aggregation of place cats / moods / pop into trip?
	moods := make([]float64, len(r.places[0].Mood))
	for i := 0; i < len(moods); i++ {
		for j := 0; j < len(r.places); j++ {
			moods[i] += r.places[j].Mood[i]
		}
		moods[i] = moods[i] / math.Sqrt(float64(len(r.places)))
	}
	return moods
}

// GetPop : Get Pop for route
func (r *Route) GetPop() float64 { // how to improve aggregation of place cats / moods / pop into trip?
	pop := 0.0
	for j := 0; j < len(r.places); j++ {
		pop += r.places[j].Pop
	}
	pop = pop / float64(len(r.places))
	return pop
}

// SetPlace : Set position of Place in Route slice
func (r *Route) SetPlace(RoutePosition int, c Place) {
	r.places[RoutePosition] = c
	// Reset fitness if Route have been altered
	r.fitness = 0
	r.distance = 0
}

func (r *Route) ResetFitnessDistance() {
	r.fitness = 0
	r.distance = 0
}

func (r *Route) RouteSize() int {
	return len(r.places)
}

// RouteDistance : Calculates total distance traveled for this Route
func (a *Route) calculateDistance() float64 {
	// if a.distance == 0 {
	// 	td := float64(0)
	// 	for i := 0; i < a.RouteSize()-1; i++ {
	// 		dist := DistMats[a.places[0].cityID]["walk_time_matrix"][a.places[i].Order][a.places[i+1].Order]
	// 		if !math.IsNaN(dist) {
	// 			td += dist
	// 		} else {
	// 			td += math.Inf(1)
	// 		}
	// 	}
	// 	a.distance = td
	// }
	if a.time == 0 {
		dists := make([]float64, a.RouteSize()-1)
		times := make([]float64, a.RouteSize()-1)
		ddists := make([]float64, a.RouteSize()-1)
		dtimes := make([]float64, a.RouteSize()-1)
		for i := 0; i < a.RouteSize()-1; i++ {
			if a.places[0].CityID[len(a.places[0].CityID)-2:] == "KR" {
				dists[i] = 111 * (math.Abs(a.places[i+1].Coords["lat"]) - a.places[i].Coords["lat"] + 88*math.Abs(a.places[i+1].Coords["lng"]-a.places[i].Coords["lng"]))
				times[i] = dists[i] / 5.5
				ddists[i] = dists[i]
				dtimes[i] = dists[i] / 50
			}
			dists[i] = DistMats[a.places[0].CityID]["walk_dist_matrix"][a.places[i].Order][a.places[i+1].Order]
			times[i] = DistMats[a.places[0].CityID]["walk_time_matrix"][a.places[i].Order][a.places[i+1].Order]
			ddists[i] = DistMats[a.places[0].CityID]["drive_dist_matrix"][a.places[i].Order][a.places[i+1].Order]
			dtimes[i] = DistMats[a.places[0].CityID]["drive_time_matrix"][a.places[i].Order][a.places[i+1].Order]
		}
		wt := float64(0)
		wd := float64(0)
		td := float64(0)
		tt := float64(0)
		// is this going to cause issues in indexing later (no)
		sort.SliceStable(times, func(i, j int) bool { return times[i] < times[j] })
		sort.SliceStable(dists, func(i, j int) bool { return times[i] < times[j] })
		sort.SliceStable(dtimes, func(i, j int) bool { return times[i] < times[j] })
		sort.SliceStable(ddists, func(i, j int) bool { return times[i] < times[j] })
		for i := 0; i < a.RouteSize()-1; i++ {
			if wd+dists[i] > 12+5*(rand.Float64()-0.5) { // driving
				if !math.IsNaN(times[i]) {
					td += ddists[i]
					tt += dtimes[i]
				} else {
					wt += math.Inf(1)
					wd += math.Inf(1)
					td += math.Inf(1)
					tt += math.Inf(1)
				}
				continue
			}
			if !math.IsNaN(times[i]) { // walking
				wt += times[i]
				wd += dists[i]
				td += dists[i]
				tt += times[i]
			} else {
				wt += math.Inf(1)
				wd += math.Inf(1)
				td += math.Inf(1)
				tt += math.Inf(1)
			}
		}
		a.walkTime = wt
		a.walkDist = wd
		a.distance = td
		a.time = tt
	}
	return a.time
}

func (a *Route) GetDistance() float64 {
	return a.distance
}

func (a *Route) GetWTime() float64 {
	return a.walkTime
}

func (a *Route) GetTime() float64 {
	return a.time
}

func (a *Route) Fitness() float64 {
	if a.fitness == 0 {
		a.calculateDistance()
		a.fitness = 1 / a.GetTime()
	}
	return a.fitness
}

func (a *Route) ContainPlace(c Place) bool {
	for _, cs := range a.places {
		if cs.PlaceID == c.PlaceID {
			return true
		}
	}
	return false
}
