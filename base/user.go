package base

import (
	"math/rand"
)

// User contains both user data and recommendations
type User struct {
	uid                   string
	name                  string
	cat                   []float64
	mood                  []float64
	pop                   float64
	walkDist              float64
	suggestedCities       []string // add sorting places to the createTrip http pipeline
	suggestedCMatch       []float64
	suggestedTrips        []string
	suggestedTMatch       []float64
	suggestedManualTrips  []string
	suggestedManualTMatch []float64
	createdTrips          []string // http pipeline
	yourTrips             []string
	yourTMatch            []float64
}

// CreateUser Returns a user struct
// firestore update, so not real-time, move real-time to a separate http-triggered method
func CreateUser(uid string, name string, lookedCities []string, likedCities []string, savedCities []string,
	lookedTrips []string, lookedTripCities []string, likedTrips []string, likedTripCities []string, savedTrips []string, savedTripCities []string,
	lookedPlaces []string, likedPlaces []string, savedPlaces []string,
	walkDist float64, newUser bool) User {
	u := User{}
	u.uid = uid
	u.name = name
	u.cat = make([]float64, len(Cities["SEATTLE_WA_US"].cat)) // any city would do
	u.mood = make([]float64, len(Cities["SEATTLE_WA_US"].mood))
	u.pop = 0
	// think about collaborative filtering
	u.walkDist = walkDist
	u.GetUserParams(lookedCities, likedCities, savedCities, lookedTrips, lookedTripCities,
		likedTrips, likedTripCities, savedTrips, savedTripCities, lookedPlaces, likedPlaces, savedPlaces, newUser)
	return u
}

// GetUserParams Put the main stuff here
func (u User) GetUserParams(lookedCities []string, likedCities []string,
	savedCities []string, lookedTrips []string, lookedTripCities []string,
	likedTrips []string, likedTripCities []string, savedTrips []string, savedTripCities []string,
	lookedPlaces []string, likedPlaces []string, savedPlaces []string, newUser bool) {
	// set user params
	weightPlaces := []float64{0.3, 0.3, 0.4}  // city, trips, place
	weightTypes := []float64{0.05, 0.45, 0.5} // looked, liked, saved
	if newUser {
		u.updateParamsNP([]float64{1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5)}, []float64{1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5), 1 + (rand.Float64()*0.5 - 0.5)}, 1)
	}
	// currently don't care about dates
	for i := 0; i < len(lookedCities); i++ { // no pop
		u.updateParamsNP(Cities[lookedCities[i]].cat, Cities[lookedCities[i]].mood,
			(1/float64(len(lookedCities)))*weightPlaces[0]*weightTypes[0])
	}
	for i := 0; i < len(likedCities); i++ {
		u.updateParamsNP(Cities[likedCities[i]].cat, Cities[likedCities[i]].mood,
			(1/float64(len(likedCities)))*weightPlaces[0]*weightTypes[1])
	}
	for i := 0; i < len(savedCities); i++ {
		u.updateParamsNP(Cities[savedCities[i]].cat, Cities[savedCities[i]].mood,
			(1/float64(len(savedCities)))*weightPlaces[0]*weightTypes[2])
	}
	for i := 0; i < len(lookedTrips); i++ {
		cityID := lookedTripCities[i]
		u.updateParams(Trips[cityID][lookedTrips[i]].Cat, Trips[cityID][lookedTrips[i]].Mood,
			Trips[cityID][lookedTrips[i]].Pop, (1/float64(len(lookedTrips)))*weightPlaces[1]*weightTypes[0])
	}
	for i := 0; i < len(likedTrips); i++ {
		cityID := likedTripCities[i]
		u.updateParams(Trips[cityID][likedTrips[i]].Cat, Trips[cityID][likedTrips[i]].Mood,
			Trips[cityID][likedTrips[i]].Pop, (1/float64(len(likedTrips)))*weightPlaces[1]*weightTypes[1])
	}
	for i := 0; i < len(savedTrips); i++ {
		cityID := savedTripCities[i]
		u.updateParams(Trips[cityID][savedTrips[i]].Cat, Trips[cityID][savedTrips[i]].Mood,
			Trips[cityID][savedTrips[i]].Pop, (1/float64(len(savedTrips)))*weightPlaces[1]*weightTypes[2])
	}
	for i := 0; i < len(lookedPlaces); i++ {
		u.updateParams(Places[PlacesCities[lookedPlaces[i]]][lookedPlaces[i]].Cat, Places[PlacesCities[lookedPlaces[i]]][lookedPlaces[i]].Mood,
			Places[PlacesCities[lookedPlaces[i]]][lookedPlaces[i]].Pop, (1/float64(len(lookedPlaces)))*weightPlaces[2]*weightTypes[0])
	}
	for i := 0; i < len(likedPlaces); i++ {
		u.updateParams(Places[PlacesCities[likedPlaces[i]]][likedPlaces[i]].Cat, Places[PlacesCities[likedPlaces[i]]][likedPlaces[i]].Mood,
			Places[PlacesCities[likedPlaces[i]]][likedPlaces[i]].Pop, (1/float64(len(likedPlaces)))*weightPlaces[2]*weightTypes[1])
	}
	for i := 0; i < len(savedPlaces); i++ {
		u.updateParams(Places[PlacesCities[savedPlaces[i]]][savedPlaces[i]].Cat, Places[PlacesCities[savedPlaces[i]]][savedPlaces[i]].Mood,
			Places[PlacesCities[savedPlaces[i]]][savedPlaces[i]].Pop, (1/float64(len(savedPlaces)))*weightPlaces[2]*weightTypes[2])
	}
	return
}

// update params given a liked set and weight coefficient
func (u User) updateParams(cat []float64, mood []float64, pop float64, weight float64) { // weight includes both the number of (places), type and age
	for i := 0; i < len(cat); i++ {
		u.cat[i] += weight * cat[i]
	}
	for i := 0; i < len(mood); i++ {
		u.mood[i] += weight * mood[i]
	}
	u.pop += weight * pop
	// dwalkdist
	return
}

// update params given a liked set and weight coefficient, but no pop
func (u User) updateParamsNP(cat []float64, mood []float64, weight float64) { // weight includes both the number of (places), type and age
	for i := 0; i < len(cat); i++ {
		u.cat[i] += weight * cat[i]
	}
	for i := 0; i < len(mood); i++ {
		u.mood[i] += weight * mood[i]
	}
	// dwalkdist
	return
}

// GetUserSuggested Put the main stuff here
func (u User) GetUserSuggested(savedTrips []string, newUser bool) (map[string]Trip, User, []City) { // add collab filtering
	rankedCities := RankCities(u.cat, u.mood) // set to fs
	for i := 0; i < 60; i++ {
		if i%3 == 0 {
			u.suggestedCities = append(u.suggestedCities, rankedCities[i].cityID)
			u.suggestedCMatch = append(u.suggestedCMatch, rankedCities[i].mc)
		} else if i%3 == 1 {
			bestTrip := RankTrips(u.cat, u.mood, u.pop, rankedCities[i])
			u.suggestedTrips = append(u.suggestedTrips, bestTrip.TripID)
			u.suggestedTMatch = append(u.suggestedTMatch, bestTrip.MC)
		} else {
			bestTrip := RankManualTrips(u.cat, u.mood, u.pop, rankedCities[i])
			if bestTrip.TripID != "" {
				u.suggestedManualTrips = append(u.suggestedManualTrips, bestTrip.TripID)
				u.suggestedManualTMatch = append(u.suggestedManualTMatch, bestTrip.MC)
			}
		}
	}
	perms := rand.Perm(10)
	for i := 0; i < 10; i++ {
		u.suggestedCities = append(u.suggestedCities, rankedCities[60+perms[i]].cityID)
		u.suggestedCMatch = append(u.suggestedCMatch, rankedCities[60+perms[i]].mc)
	}
	// returns a map, add not only the trip id, mc to user, but also your trips
	// to trips, also delete the old sometimes
	yourTrips := make(map[string]Trip)
	if rand.Float64() < 0.3 && !newUser {
		// pick random 3 from top 15, replace, unless saved, remove ones that are not in top 10, unless saved
		cities := make([]City, 0)
		perms := rand.Perm(15)
		for _, idx := range []int{0, 1, 2} {
			// exists := false
			// for _, city := range savedTrips {
			// 	if rankedCities[perms[idx]].GetCityID()+"_U_"+u.GetID() == city {
			// 		exists = true
			// 	}
			// }
			// if !exists {
			// 	cities = append(cities, rankedCities[perms[idx]])
			// }
			cities = append(cities, rankedCities[perms[idx]])
		}
		yourTrips = MakeYourTrips(cities, u.cat, u.mood, u.pop, u.uid) // for you (not created -> http)
	}
	for _, trip := range yourTrips {
		u.yourTrips = append(u.yourTrips, trip.CityID)
		u.yourTMatch = append(u.yourTMatch, trip.MC)
	}
	// TODO: Suggest others' trips
	return yourTrips, u, rankedCities[:15]
}

// GetID gets User ID
func (u User) GetID() string {
	return u.uid
}

// GetCats gets cats
func (u User) GetCats() []float64 {
	return u.cat
}

// GetMoods gets moods
func (u User) GetMoods() []float64 {
	return u.mood
}

// GetPop gets pop
func (u User) GetPop() float64 {
	return u.pop
}

// GetSuggestedCities gets suggested cities
func (u User) GetSuggestedCities() []string {
	return u.suggestedCities
}

// GetSuggestedTrips gets suggested cities
func (u User) GetSuggestedTrips() []string {
	return u.suggestedTrips
}

// GetSuggestedManualTrips gets suggested cities (created)
func (u User) GetSuggestedManualTrips() []string {
	return u.suggestedManualTrips
}

// GetName gets pop
func (u User) GetName() string {
	return u.name
}
