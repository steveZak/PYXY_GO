package base

import (
	"math"
	"math/rand"
	"sync"

	"github.com/rs/xid"
	"gonum.org/v1/gonum/stat/combin"
)

// Trip : trip struct
type Trip struct {
	UserID   string    `json:"userID"`
	CityID   string    `json:"cityID"`
	TripID   string    `json:"tripID"`
	Source   string    `json:"source"`
	Tag      string    `json:"tag"`
	Places   []Place   `json:"places"`
	MC       float64   `json:"mc"`
	Cat      []float64 `json:"cat"`
	Mood     []float64 `json:"mood"`
	Privacy  string    `json:"privacy"`
	Pop      float64   `json:"pop"`
	Time     float64   `json:"time"`
	WalkDist float64   `json:"walkDist"`
	Distance float64   `json:"distance"`
}

// GenerateTrip : generate trip and set params (setting MC to -1 for trips with Source 'auto')
func GenerateTrip(UserID string, CityID string, TripID string, Source string, Tag string, Places []Place, Cat []float64, Mood []float64, Pop float64, Privacy string,
	Time float64, WalkDist float64, Distance float64, MC float64) Trip {
	t := Trip{}
	t.UserID = UserID
	t.CityID = CityID
	t.TripID = TripID
	t.Source = Source
	t.Tag = Tag
	t.Places = Places
	t.Cat = Cat
	t.Mood = Mood
	t.Pop = Pop
	t.Privacy = Privacy
	t.Time = Time
	t.WalkDist = WalkDist
	t.Distance = Distance
	t.MC = MC
	return t
}

// ImportTrip : import trip and set params (setting MC to -1 for trips with Source 'auto')
func ImportTrip(CityID string, TripID string, Source string, Tag string, Places []Place, Cat []interface{}, Mood []interface{}, Pop float64,
	Time float64, WalkDist float64, Distance float64, MC float64) Trip {
	t := Trip{}
	t.CityID = CityID
	t.TripID = TripID
	t.Source = Source
	t.Tag = Tag
	t.Places = Places
	for _, v := range Cat {
		t.Cat = append(t.Cat, v.(float64))
	}
	for _, v := range Mood {
		t.Mood = append(t.Mood, v.(float64))
	}
	t.Pop = Pop
	t.Privacy = "public"
	t.Time = Time
	t.WalkDist = WalkDist
	t.Distance = Distance
	t.MC = MC
	return t
}

// DownloadTrip : download trip and set params (setting MC to -1 for trips with Source 'auto')
func DownloadTrip(UserID string, CityID string, TripID string, Source string, Tag string, Places []Place, Cat []interface{}, Mood []interface{}, Pop float64, Privacy string,
	Time float64, WalkDist float64, Distance float64, MC float64) Trip {
	t := Trip{}
	t.UserID = UserID
	t.CityID = CityID
	t.TripID = TripID
	t.Source = Source
	t.Tag = Tag
	t.Places = Places
	for _, v := range Cat {
		t.Cat = append(t.Cat, v.(float64))
	}
	for _, v := range Mood {
		t.Mood = append(t.Mood, v.(float64))
	}
	t.Pop = Pop
	t.Privacy = Privacy
	t.Time = Time
	t.WalkDist = WalkDist
	t.Distance = Distance
	t.MC = MC
	return t
}

// cosine similarity
func (t *Trip) GetTripMC(Cat []float64, Mood []float64, Pop float64) float64 {
	sumUC, sumUU, suMCC := 0.0, 0.0, 0.0
	for i := 0; i < len(Cat); i++ {
		sumUC += Cat[i] * t.Cat[i]
		sumUU += Cat[i] * Cat[i]
		suMCC += t.Cat[i] * t.Cat[i]
	}
	kC := sumUC / (math.Sqrt(sumUU) * math.Sqrt(suMCC))
	sumUC, sumUU, suMCC = 0.0, 0.0, 0.0
	for i := 0; i < len(Mood); i++ {
		sumUC += Mood[i] * t.Mood[i]
		sumUU += Mood[i] * Mood[i]
		suMCC += t.Mood[i] * t.Mood[i]
	}
	kM := sumUC / (math.Sqrt(sumUU) * math.Sqrt(suMCC))
	kP := 1 - math.Abs(Pop-t.Pop)
	return 0.75*(kC+1.0)/2.0 + 0.15*(kM+1.0)/2.0 + 0.1*kP
}

// // RankTrips : not sorting by Time, because unable to vary trips
// func RankTrips(Cat []float64, Mood []float64, Pop float64) []Trip { //consider: match countries
// 	var t []Trip
// 	for _, cityTrips := range Trips {
// 		for _, trip := range cityTrips {
// 			t = append(t, trip)
// 		}
// 	}
// 	t[0].MC = t[0].GetTripMC(Cat, Mood, Pop) // change to highest cats/moods trips for a city
// 	tSort := []Trip{t[0]}
// 	for i := 1; i < len(t); i++ {
// 		t[i].MC = t[i].GetTripMC(Cat, Mood, Pop)
// 		j := sort.Search(len(tSort), func(j int) bool { return tSort[j].MC > t[i].MC })
// 		if j == len(tSort) {
// 			tSort = append([]Trip{t[i]}, tSort...)
// 			continue
// 		}
// 		if j == len(tSort)-1 {
// 			tSort = append(tSort, []Trip{t[i]}...)
// 			continue
// 		}
// 		tSort = append(tSort[0:j], append([]Trip{t[i]}, tSort[j+1:]...)...)
// 	}
// 	return tSort
// }

// RankTrips : not sorting by Time, because unable to vary trips
func RankTrips(Cat []float64, Mood []float64, Pop float64, city City) Trip { //consider: match countries
	var t []Trip
	for _, trip := range Trips[city.GetCityID()] {
		if trip.Tag == "generated" {
			t = append(t, trip)
		}
	}
	var bt Trip
	maxmc := 0.0
	for i := 0; i < len(t); i++ {
		t[i].MC = t[i].GetTripMC(Cat, Mood, Pop)
		if t[i].MC > maxmc {
			maxmc = t[i].MC
			bt = t[i]
		}
	}
	return bt
}

// RankTrips : not sorting by Time, because unable to vary trips
func RankManualTrips(Cat []float64, Mood []float64, Pop float64, city City) Trip { //consider: match countries
	var t []Trip
	for _, trip := range Trips[city.GetCityID()] {
		if trip.Tag == "created" && trip.Privacy == "public" {
			t = append(t, trip)
		}
	}
	var bt Trip
	if len(t) == 0 {
		return bt
	}
	maxmc := 0.0
	for i := 0; i < len(t); i++ {
		t[i].MC = t[i].GetTripMC(Cat, Mood, Pop)
		if t[i].MC > maxmc {
			maxmc = t[i].MC
			bt = t[i]
		}
	}
	return bt
}

// MakeYourTrips : create trips for the user-specific params for top fitting cities
func MakeYourTrips(c []City, Cat []float64, Mood []float64, Pop float64, UID string) map[string]Trip { //consider: match countries
	yourTrips := make(map[string]Trip)
	for _, city := range c {
		places := Places[city.GetCityID()]
		var ranked []Place
		for _, place := range places { // untested
			ranked = append(ranked, place)
		}
		ranked = RankPlaces(ranked, Cat, Mood, Pop, 10)
		// run GA (generateTrip, like in import)
		c := make(chan Route) // verify that this works
		var wg sync.WaitGroup
		num := 0
		for i := 5; i < 9; i++ {
			combs := combin.Combinations(8, i)
			for j := 0; j < len(combs); j++ {
				tm := TripManager{}
				tm.NewTripManager(10)
				for k := range combs[j] {
					tm.AddPlace(ranked[combs[j][k]])
				}
				wg.Add(1)
				go tspGA(&wg, &tm, c, 25, num) //12
				num++
			}
		}
		routes := make([]Route, num)
		mcs := make([]float64, num)
		for j := 0; j < num; j++ {
			route := <-c
			route.SetMatch(Cat, Mood, Pop)
			routes[j] = route
			mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[city.GetCityID()]["avg_dd"][0][0]*150.0)
		}
		wg.Wait()
		maxIdx := 0 // get the most efficient trips
		maxMC := 0.0
		for idx, e := range mcs {
			if idx == 0 || maxMC < e {
				maxMC = e
				maxIdx = idx
			}
		}
		t := GenerateTrip(UID, city.GetCityID(), xid.New().String(), "auto", "your", routes[maxIdx].GetPlaces(),
			routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "private",
			routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
		yourTrips[t.TripID] = t
		Trips[t.CityID][t.TripID] = t
	}
	return yourTrips
}

// MakeYourTrip : make your trip in a user-selected city (http/firestore pipeline)
func MakeYourTrip(cityID string, tripID string, Cat []float64, Mood []float64, Pop float64, UID string) Trip { //consider: match countries
	var t Trip
	places := Places[cityID]
	var ranked []Place
	for _, place := range places { // untested
		ranked = append(ranked, place)
	}
	ranked = RankPlaces(ranked, Cat, Mood, Pop, 11)
	// run GA (generateTrip, like in import)
	c := make(chan Route) // verify that this works
	var wg sync.WaitGroup
	num := 0
	for i := 5; i < 9; i++ {
		combs := combin.Combinations(8, i)
		for j := 0; j < len(combs); j++ {
			tm := TripManager{}
			tm.NewTripManager(10)
			for k := range combs[j] {
				tm.AddPlace(ranked[combs[j][k]])
			}
			wg.Add(1)
			go tspGA(&wg, &tm, c, 25, num)
			num++
		}
	}
	routes := make([]Route, num)
	mcs := make([]float64, num)
	for j := 0; j < num; j++ {
		route := <-c
		route.SetMatch(Cat, Mood, Pop)
		routes[j] = route
		mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[cityID]["avg_dd"][0][0]*150.0) //160
	}
	wg.Wait()
	maxIdx := 0 // get the most efficient trips
	maxMC := 0.0
	for idx, e := range mcs {
		if idx == 0 || maxMC < e {
			maxMC = e
			maxIdx = idx
		}
	}
	t = GenerateTrip(UID, cityID, tripID, "auto", "your", routes[maxIdx].GetPlaces(),
		routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "private",
		routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
	Trips[t.CityID][t.TripID] = t
	return t
}

// CreateTrip : create the trip after the user defined all places (http pipeline)
func CreateTrip(cityID string, tripID string, places []Place, UID string, privacy string) Trip { //consider: match countries
	// run GA (generateTrip, like in import)
	tm := TripManager{}
	tm.NewTripManager(len(places))
	for _, place := range places {
		tm.AddPlace(place)
	}
	route := tspSingleGA(&tm, 25)
	t := GenerateTrip(UID, cityID, tripID, "manual", "created", route.GetPlaces(),
		route.GetCats(), route.GetMoods(), route.GetPop(), privacy,
		route.time, route.walkDist, route.distance, -1)
	Trips[t.CityID][t.TripID] = t
	return t
}

// EditTrip : edit the trip to change privacy(http pipeline)
func EditTrip(cityID string, tripID string, newTripID string, places []Place, UID string, privacy string) Trip { //consider: match countries
	// do not run GA (places order predetermined, just make a route to determine a few fields)
	tm := TripManager{}
	tm.NewTripManager(len(places))
	for _, place := range places {
		tm.AddPlace(place)
	}
	r := Route{}
	r.InitRoutePlaces(tm)
	r.Fitness()
	t := Trips[cityID][tripID]
	nt := GenerateTrip(UID, cityID, tripID, t.Source, t.Tag, places,
		r.GetCats(), r.GetMoods(), r.GetPop(), privacy,
		r.time, r.walkTime, r.distance, -1)
	Trips[cityID][newTripID] = nt
	return nt
}

// ShuffleTrips : if params are all default
func ShuffleTrips(in []Trip) []Trip {
	out := make([]Trip, len(in), cap(in))
	perm := rand.Perm(len(in))
	for i, v := range perm {
		out[v] = in[i]
	}
	return out
}
