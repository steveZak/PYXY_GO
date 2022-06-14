package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"gonum.org/v1/gonum/stat/combin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	// Cities cities
	Cities = importCities()
	// Places places
	Places = importPlaces()
	// PlacesCities dict of places to cities
	PlacesCities = processPlaces()
	// Users users
	Users = getUsers()
	// DistMats distance matrices
	DistMats = importDistMats()
	// Trips trips
	Trips = importTrips()
	// Trips = generateTrips()
)

func importCities() map[string]City {
	fileJ, _ := os.Open("data/cities.json")
	data, _ := ioutil.ReadAll(fileJ)
	var citiesData map[string]interface{}
	json.Unmarshal([]byte(data), &citiesData)
	var cities = make(map[string]City, len(citiesData))
	for k := range citiesData { // iterates over cities
		cities[k] = GenerateCity(k,
			citiesData[k].(map[string]interface{})["cat_params"].([]interface{}),
			citiesData[k].(map[string]interface{})["mood_params"].([]interface{}))
	}
	return cities
}

func importPlaces() map[string]map[string]Place {
	fileJ, _ := os.Open("data/places.json")
	data, _ := ioutil.ReadAll(fileJ)
	var placesData map[string]map[string]interface{}
	var places = make(map[string]map[string]Place)
	json.Unmarshal([]byte(data), &placesData)
	for k, city := range placesData { // iterates over cities
		sights := city["sights"].([]interface{})
		var placesCity = make(map[string]Place)
		for i, place := range sights {
			placesCity[place.(map[string]interface{})["place_id"].(string)] = GeneratePlace(k,
				place.(map[string]interface{})["place_id"].(string),
				i,
				place.(map[string]interface{})["name"].(string),
				place.(map[string]interface{})["global_cat_params"].([]interface{}),
				place.(map[string]interface{})["global_mood_params"].([]interface{}),
				place.(map[string]interface{})["popularity"].(float64),
				place.(map[string]interface{})["coordinates"].(map[string]interface{}),
				place.(map[string]interface{})["duration"].(float64))
		}
		places[k] = placesCity
	}
	return places
}

func processPlaces() map[string]string {
	var places = make(map[string]string)
	for _, cityPlaces := range Places {
		for _, place := range cityPlaces {
			places[place.PlaceID] = place.CityID
		}
	}
	return places
}

func getUsers() map[string]User {
	var users = make(map[string]User)
	// opt := option.WithCredentialsFile("cloud_access/firebase_credentials.json")
	opt := option.WithCredentialsJSON([]byte(`{
		"type": "***",
		"project_id": "***",
		"private_key_id": "6***c67c413d7a728307c0e44afbd38e69a07998a09",
		"private_key": "***",
		"client_email": "***",
		"client_id": "***",
		"auth_uri": "***",
		"token_uri": "***",
		"auth_provider_x509_cert_url": "***",
		"client_x509_cert_url": ""
	  }`))
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "***", opt)
	if err != nil {
		// TODO: Handle error.
		fmt.Println(err)
	}
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error
		}
		data := doc.Data()
		lookedCities := make([]string, 0)
		if data["looked_cities"] != nil {
			for _, v := range data["looked_cities"].([]interface{}) {
				vv := v.(map[string]interface{})
				lookedCities = append(lookedCities, vv["city_id"].(string))
				// lookedCDates = append(lookedCDates, vv["date"].(string))
			}
		}
		likedCities := make([]string, 0)
		if data["liked_cities"] != nil {
			for _, v := range data["liked_cities"].([]interface{}) {
				// vv := v.(map[string]interface{})
				// likedCities = append(likedCities, vv["city_id"].(string))
				vv := v.(string)
				likedCities = append(likedCities, vv)
			}
		}
		savedCities := make([]string, 0)
		if data["saved_cities"] != nil {
			for _, v := range data["saved_cities"].([]interface{}) {
				// vv := v.(map[string]interface{})
				// savedCities = append(savedCities, vv["city_id"].(string))
				vv := v.(string)
				savedCities = append(savedCities, vv)
			}
		}
		lookedPlaces := make([]string, 0)
		if data["looked_places"] != nil { // remember to switch 'sights' to 'places' in firestore
			for _, v := range data["looked_places"].([]interface{}) {
				vv := v.(map[string]interface{})
				lookedPlaces = append(lookedPlaces, vv["place_id"].(string))
				// lookedPDates = append(lookedPDates, vv["date"].(string))
			}
		}
		likedPlaces := make([]string, 0)
		if data["liked_places"] != nil {
			for _, v := range data["liked_places"].([]interface{}) {
				// vv := v.(map[string]interface{})
				// likedPlaces = append(likedPlaces, vv["city_id"].(string))
				vv := v.(string)
				likedPlaces = append(likedPlaces, vv)
			}
		}
		savedPlaces := make([]string, 0)
		if data["saved_places"] != nil {
			for _, v := range data["saved_places"].([]interface{}) {
				// vv := v.(map[string]interface{})
				// savedPlaces = append(savedPlaces, vv["place_id"].(string))
				vv := v.(string)
				savedPlaces = append(savedPlaces, vv)
			}
		}
		lookedTrips := make([]string, 0)
		lookedTripCities := make([]string, 0)
		if data["looked_trips"] != nil {
			for _, v := range data["looked_trips"].([]interface{}) {
				vv := v.(map[string]interface{})
				lookedTrips = append(lookedTrips, vv["trip_id"].(string))
				lookedTripCities = append(lookedTripCities, vv["city_id"].(string))
			}
		}
		likedTrips := make([]string, 0)
		likedTripCities := make([]string, 0)
		if data["liked_trips"] != nil {
			for _, v := range data["liked_trips"].([]interface{}) {
				vv := v.(map[string]interface{})
				likedTrips = append(likedTrips, vv["trip_id"].(string))
				likedTripCities = append(likedTripCities, vv["city_id"].(string))
			}
		}
		savedTrips := make([]string, 0)
		savedTripCities := make([]string, 0)
		if data["saved_trips"] != nil {
			for _, v := range data["saved_trips"].([]interface{}) {
				vv := v.(map[string]interface{})
				savedTrips = append(savedTrips, vv["trip_id"].(string))
				savedTripCities = append(savedTripCities, vv["city_id"].(string))
			}
		}
		newUser := len(lookedCities)+len(likedCities)+len(savedCities)+len(lookedTrips)+len(likedTrips)+len(savedTrips)+len(lookedPlaces)+len(likedPlaces)+len(savedPlaces) == 0
		users[doc.Ref.ID] = CreateUser(doc.Ref.ID, data["name"].(string), lookedCities, likedCities, savedCities,
			lookedTrips, lookedTripCities, likedTrips, likedTripCities, savedTrips,
			savedTripCities, lookedPlaces, likedPlaces, savedPlaces, float64(data["walk_dist"].(int64)), newUser)
	}
	return users
}

func importDistMats() map[string]map[string][][]float64 {
	fileJ, _ := os.Open("data/dist_mat.json")
	data, _ := ioutil.ReadAll(fileJ)
	data = bytes.Replace(data, []byte("NaN"), []byte("null"), -1)
	var distMats map[string]map[string][][]float64
	json.Unmarshal([]byte(data), &distMats)
	return distMats
}

func importTrips() map[string]map[string]Trip { // lacks coords & duration
	fileJ, _ := os.Open("data/trips.json")
	data, _ := ioutil.ReadAll(fileJ)
	var allTripsData map[string]map[string]interface{}
	allTrips := make(map[string]map[string]Trip)
	json.Unmarshal([]byte(data), &allTripsData)
	for cityID, cityTrips := range allTripsData { // iterates over cities
		var trips = make(map[string]Trip)
		for tripID, trip := range cityTrips {
			sights := trip.(map[string]interface{})["places"].([]interface{})
			trips[tripID] = ImportTrip(trip.(map[string]interface{})["cityID"].(string),
				trip.(map[string]interface{})["tripID"].(string),
				"auto",
				"generated",                // trip.(map[string]interface{})["tag"].(string)
				make([]Place, len(sights)), // make an array of places for this trip
				trip.(map[string]interface{})["cat"].([]interface{}),
				trip.(map[string]interface{})["mood"].([]interface{}),
				trip.(map[string]interface{})["pop"].(float64),
				trip.(map[string]interface{})["time"].(float64),
				trip.(map[string]interface{})["walkDist"].(float64),
				trip.(map[string]interface{})["distance"].(float64),
				-1)
			for i, place := range sights { // loop over the places in a trip
				trips[tripID].Places[i] = ImportPlace(place.(map[string]interface{})["cityID"].(string),
					place.(map[string]interface{})["placeID"].(string),
					place.(map[string]interface{})["name"].(string),
					place.(map[string]interface{})["cat"].([]interface{}),
					place.(map[string]interface{})["mood"].([]interface{}),
					place.(map[string]interface{})["pop"].(float64))
			}
		}
		allTrips[cityID] = trips
	}
	// download trips as well
	opt := option.WithCredentialsJSON([]byte(`{
		"type": "***",
		"project_id": "***",
		"private_key_id": "***",
		"private_key": "***",
		"client_email": "***",
		"client_id": "***",
		"auth_uri": "***",
		"token_uri": "***",
		"auth_provider_x509_cert_url": "***",
		"client_x509_cert_url": "***"
	  }`))
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "***", opt)
	if err != nil {
		fmt.Println(err)
	}
	iter := client.Collection("trips").Where("tag", "==", "created").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error
		}
		data := doc.Data()
		places := make([]Place, 0)
		if data["sights"] != nil {
			for _, v := range data["sights"].([]interface{}) {
				vv := v.(map[string]interface{})["place_id"].(string)
				places = append(places, Places[data["city_id"].(string)][vv])
			}
		}
		allTrips[data["city_id"].(string)][doc.Ref.ID] = DownloadTrip(data["user_id"].(string), data["city_id"].(string), doc.Ref.ID, "manual",
			"created", places, data["global_cat_params"].([]interface{}), data["global_mood_params"].([]interface{}), data["pop"].(float64), data["privacy"].(string),
			data["time"].(float64), data["walk_dist"].(float64), data["distance"].(float64), -1)
	}
	iter = client.Collection("trips").Where("tag", "==", "yours").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error
		}
		data := doc.Data()
		places := make([]Place, 0)
		if data["sights"] != nil {
			for _, v := range data["sights"].([]interface{}) {
				vv := v.(map[string]interface{})["place_id"].(string)
				places = append(places, Places[data["city_id"].(string)][vv])
			}
		}
		allTrips[data["city_id"].(string)][doc.Ref.ID] = ImportTrip(data["city_id"].(string), doc.Ref.ID, "auto",
			"yours", places, data["global_cat_params"].([]interface{}), data["global_mood_params"].([]interface{}), data["pop"].(float64),
			data["time"].(float64), data["walk_dist"].(float64), data["distance"].(float64), -1)
	}
	// iter = client.Collection("trips").Where("tag", "==", "created").Documents(ctx)
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		// Handle error
	// 	}
	// 	data := doc.Data()
	// 	places := make([]Place, 0)
	// 	if data["places"] != nil {
	// 		for _, v := range data["places"].([]interface{}) {
	// 			vv := v.(map[string]interface{})["place_id"].(string)
	// 			places = append(places, Places[data["city_id"].(string)][vv])
	// 		}
	// 	}
	// 	allTrips[data["city_id"].(string)][data["user_id"].(string)] = ImportTrip(data["city_id"].(string), doc.Ref.ID, "auto",
	// 		"created", places, data["global_cat_params"].([]interface{}), data["global_mood_params"].([]interface{}), data["pop"].(float64),
	// 		data["time"].(float64), data["walk_dist"].(float64), data["distance"].(float64), -1)
	// }
	return allTrips
}

func generateTrips() map[string]map[string]Trip {
	allTrips := make(map[string]map[string]Trip)
	catsLabels := []string{"Amusement", "Animals", "Architecture", "Art", "Beach", "Books", "Cosmopolitan", "Culture", "Engineering", "Family", "Fashion", "Food", "Friends", "Hiking", "History", "Hospitality", "Insta", "Military", "Music", "Movies", "Nature", "Nightlife", "Original", "Religion", "Science", "Shopping", "Sport", "Views", "Walking", "Water"}
	moodsLabels := []string{"Adventuresome", "Entertaining", "Peaceful", "Captivating", "Magical", "Exciting", "Impressive", "Inspiring", "Romantic", "Scary", "Wise", "Sad"}
	cnt := 0
	for cityID, cityPlacesDict := range Places {
		if _, ok := DistMats[cityID]; !ok {
			fmt.Println(cityID + "is missing distmat")
			continue
		}
		cnt++
		fmt.Println(cityID)
		// if cityID != "SEOUL_SE_KR" {
		// 	continue
		// }
		cityPlaces := []Place{}
		for _, value := range cityPlacesDict {
			cityPlaces = append(cityPlaces, value)
		}
		trips := make(map[string]Trip)
		cats := Cities[cityID].cat
		moods := Cities[cityID].mood
		rankedCats := []int{}
		rankedMoods := []int{}
		for i := 0; i < len(cats); i++ {
			j := -1
			for k := 0; k < len(rankedCats); k++ {
				if cats[i] > cats[rankedCats[k]] {
					j = k
					break
				}
			}
			if j == -1 {
				j = len(rankedCats)
			}
			if j == len(rankedCats) {
				rankedCats = append(rankedCats, []int{i}...)
				continue
			}
			if j == 0 {
				rankedCats = append([]int{i}, rankedCats...)
				continue
			}
			rankedCats = append(rankedCats[0:j], append([]int{i}, rankedCats[j:]...)...)
		}
		for i := 0; i < len(moods); i++ {
			j := -1
			for k := 0; k < len(rankedMoods); k++ {
				if moods[i] > moods[rankedMoods[k]] {
					j = k
					break
				}
			}
			if j == -1 {
				j = len(rankedMoods)
			}
			if j == len(rankedMoods) {
				rankedMoods = append(rankedMoods, []int{i}...)
				continue
			}
			if j == 0 {
				rankedMoods = append([]int{i}, rankedMoods...)
				continue
			}
			rankedMoods = append(rankedMoods[0:j], append([]int{i}, rankedMoods[j:]...)...)
		}
		rankedCats = rankedCats[:10]
		// rank places for each selected param, optimise trip,
		for i := 0; i < len(rankedCats); i++ {
			catsVec := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			catsVec[rankedCats[i]] = 1.0
			c := make(chan Route) // verify that this works
			var wg sync.WaitGroup
			rankedFull := RankPlacesCat(cityPlaces, catsVec, 11) // reduce to 10 later for theme / mood / pop
			ranked := rankedFull[:10]
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
				route.SetMatchCat(catsVec)
				routes[j] = route
				mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[cityID]["avg_dd"][0][0]*150.0)
			}
			wg.Wait()
			maxIdx := 0 // get the most efficient trips
			maxMC := 0.0
			for idx, e := range mcs {
				if idx == 0 || maxMC < e { // ???
					maxMC = e
					maxIdx = idx
				}
			}
			fmt.Println(len(routes[maxIdx].GetPlaces()))
			trips[cityID+"_T_"+catsLabels[rankedCats[i]]] = GenerateTrip("", cityID, cityID+"_T_"+catsLabels[rankedCats[i]], "auto", "cat", routes[maxIdx].GetPlaces(),
				routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "public",
				routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
			var places []map[string]interface{}
			duration := 0.0
			for _, place := range trips[cityID+"_T_"+catsLabels[rankedCats[i]]].Places {
				places = append(places, map[string]interface{}{"coordinates": Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
				duration += Places[place.CityID][place.PlaceID].Duration
			}
			duration += trips[cityID+"_T_"+catsLabels[rankedCats[i]]].Time / 60.0
			fmt.Println(duration)
		}
		rankedMoods = rankedMoods[:5]
		for i := 0; i < len(rankedMoods); i++ {
			moodsVec := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			moodsVec[rankedMoods[i]] = 1.0
			c := make(chan Route) // verify that this works
			var wg sync.WaitGroup
			rankedFull := RankPlacesMood(cityPlaces, moodsVec, 11) // reduce to 10 later for theme / mood / pop
			ranked := rankedFull[:10]
			num := 0
			for ii := 5; ii < 9; ii++ {
				combs := combin.Combinations(8, ii)
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
				route.SetMatchMood(moodsVec)
				routes[j] = route
				mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[cityID]["avg_dd"][0][0]*150.0)
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
			// fmt.Println(len(routes[maxIdx].places))
			trips[cityID+"_M_"+moodsLabels[rankedMoods[i]]] = GenerateTrip("", cityID, cityID+"_M_"+moodsLabels[rankedMoods[i]], "auto", "mood", routes[maxIdx].GetPlaces(),
				routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "public",
				routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
		}
		c := make(chan Route) // verify that this works
		var wg sync.WaitGroup
		rankedFull := RankPlacesPop(cityPlaces, 1.0, 11) // reduce to 10 later for theme / mood / pop
		ranked := rankedFull[:10]
		num := 0
		for ii := 5; ii < 9; ii++ {
			combs := combin.Combinations(8, ii)
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
			route.SetMatchPop(1.0)
			routes[j] = route
			mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[cityID]["avg_dd"][0][0]*150.0)
		}
		wg.Wait()
		maxIdx := 0 // get the most efficient trip
		maxMC := 0.0
		for idx, e := range mcs {
			if idx == 0 || maxMC < e {
				maxMC = e
				maxIdx = idx
			}
		}
		// fmt.Println(len(routes[maxIdx].places))
		trips[cityID+"_P_Popular"] = GenerateTrip("", cityID, cityID+"_P_Popular", "auto", "pop", routes[maxIdx].GetPlaces(),
			routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "public",
			routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
		c = make(chan Route)                            // verify that this works
		rankedFull = RankPlacesPop(cityPlaces, 0.0, 11) // reduce to 10 later for theme / mood / pop
		ranked = rankedFull[:10]
		num = 0
		for ii := 5; ii < 9; ii++ {
			combs := combin.Combinations(8, ii)
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
		routes = make([]Route, num)
		mcs = make([]float64, num)
		for j := 0; j < num; j++ {
			route := <-c
			route.SetMatchPop(0.0)
			routes[j] = route
			mcs[j] = route.GetMatch() - route.GetTime()*route.GetTime()/(DistMats[cityID]["avg_dd"][0][0]*150.0)
		}
		wg.Wait()
		maxIdx = 0 // get the most efficient trip
		maxMC = 0.0
		for idx, e := range mcs {
			if idx == 0 || maxMC < e {
				maxMC = e
				maxIdx = idx
			}
		}
		// fmt.Println(len(routes[maxIdx].places))
		trips[cityID+"_P_Rare"] = GenerateTrip("", cityID, cityID+"_P_Rare", "auto", "pop", routes[maxIdx].GetPlaces(),
			routes[maxIdx].GetCats(), routes[maxIdx].GetMoods(), routes[maxIdx].GetPop(), "public",
			routes[maxIdx].time, routes[maxIdx].walkDist, routes[maxIdx].distance, -1)
		close(c)
		allTrips[cityID] = trips
	}
	data, err := json.MarshalIndent(allTrips, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("data/trips.json", data, 0644)
	return allTrips // create a script to upload to fb
}

// UploadTrips Uploads trips to fs
func UploadTrips() {
	// opt := option.WithCredentialsFile("cloud_access/firebase_credentials.json")
	opt := option.WithCredentialsJSON([]byte(`{
		"type": "***",
		"project_id": "",
		"private_key_id": "",
		"private_key": "***",
		"client_email": "***",
		"client_id": "***",
		"auth_uri": "***",
		"token_uri": "***",
		"auth_provider_x509_cert_url": "***",
		"client_x509_cert_url": "***"
	  }`))
	ctx := context.Background()
	client, _ := firestore.NewClient(ctx, "***", opt)
	batch := client.Batch()
	i := 0
	for _, city := range Trips {
		for _, trip := range city {
			i++
			var places []map[string]interface{}
			duration := 0.0
			for _, place := range trip.Places {
				places = append(places, map[string]interface{}{"coordinates": Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
				duration += Places[place.CityID][place.PlaceID].Duration
			}
			duration += trip.Time / 60.0
			tripRef := client.Collection("trips").Doc(trip.TripID)
			tripFB := map[string]interface{}{
				"city_id":            trip.CityID, //include in addition to existing trips
				"distance":           trip.Distance,
				"global_cat_params":  trip.Cat,
				"global_mood_params": trip.Mood,
				"pop":                trip.Pop,
				"sights":             places,
				"source":             "auto",
				"tag":                "generated",
				"time":               trip.Time,
				"walk_dist":          trip.WalkDist,
				"duration":           duration}
			batch.Set(tripRef, tripFB)
		}
		if i > 400 {
			_, err := batch.Commit(ctx)
			batch = client.Batch()
			fmt.Println(i)
			i = 0
			if err != nil {
				log.Printf("An error has occurred: %s", err)
			}
		}
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
}

func tspGA(wg *sync.WaitGroup, tm *TripManager, c chan Route, gen int, pt int) {
	defer wg.Done()
	p := Population{}
	p.InitPopulation(600, *tm)
	// iFit := p.GetFittest()
	// if pt == 609 {
	// 	fmt.Printf("init %d = %v\n", pt, iFit.GetTime())
	// }
	for i := 1; i < 9+1; i++ { // gen+1
		p = EvolvePopulation(p)
		// fFit := p.GetFittest()
		// if pt == 609 {
		// 	fmt.Printf("%d at %d = %v\n", pt, i, fFit.GetTime())
		// }
	}
	fFit := p.GetFittest()
	// if pt == 609 {
	// 	fmt.Printf("final %d = %v\n", pt, fFit.GetTime())
	// 	fmt.Printf("final %d = %v\n", pt, fFit.GetPlaces())
	// }
	c <- *fFit
}

func tspSingleGA(tm *TripManager, gen int) Route {
	p := Population{}
	p.InitPopulation(500, *tm)
	for i := 1; i < 9+1; i++ { // gen+1
		p = EvolvePopulation(p)
	}
	fFit := p.GetFittest()
	return *fFit
}
