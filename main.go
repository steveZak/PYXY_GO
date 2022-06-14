package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-pyxy/base"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/rs/xid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
//seed = time.Now().Unix()
)

func main() {
	// upload generated Trips to FS
	// cities & places still updated with python scripts
	// base.UploadTrips()

	// initialise the firebase client
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go launchHTTPServer()
	go launchListener()
	wg.Wait()
}

func launchListener() {
	fmt.Println("starting firestore listener")
	// opt := option.WithCredentialsFile("cloud_access/firebase_credentials.json")
	opt := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
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
		// TODO: Handle error.
		fmt.Println(err)
	}
	q := client.Collection("users")
	qsnapIter := q.Snapshots(ctx)
	for {
		qsnap, err := qsnapIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error
		}
		// fmt.Printf("At %s there were %d results.\n", qsnap.ReadTime, qsnap.Size)
		// docs := qsnap.Documents // TODO: Iterate over the results if desired.

		changes := qsnap.Changes // TODO: Use the list of incremental changes if desired.
		if len(changes) > 0 {
			for _, change := range changes {
				go ProcessChange(ctx, client, change)
			}
		}
	}
}

// ProcessChange concurrent thread processing the update
func ProcessChange(ctx context.Context, client *firestore.Client, change firestore.DocumentChange) {
	data := change.Doc.Data()
	needsUpdates := data["needs_update"].(bool)
	if needsUpdates {
		// fmt.Println("updating")
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
				lookedTripCities = append(lookedTripCities, vv["city_id"].(string)) // oopsie
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
		user := base.CreateUser(change.Doc.Ref.ID, data["name"].(string), lookedCities, likedCities, savedCities,
			lookedTrips, lookedTripCities, likedTrips, likedTripCities, savedTrips,
			savedTripCities, lookedPlaces, likedPlaces, savedPlaces, float64(data["walk_dist"].(int64)), newUser)
		yourTrips, user, ranked := user.GetUserSuggested(savedTrips, newUser)
		UpdateFirestore(ctx, user, yourTrips, ranked, client)
	}
}

// UpdateFirestore update specific user fields and yourTrips in batch firestore (fs pipeline)
func UpdateFirestore(ctx context.Context, u base.User, yourTrips map[string]base.Trip, ranked []base.City, client *firestore.Client) {
	base.Users[u.GetID()] = u
	user := client.Collection("users").Doc(u.GetID())
	batch := client.Batch()
	// var sc []map[string]interface{}
	// for _, city := range u.GetSuggestedCities() {
	// 	sc = append(sc, map[string]interface{}{"city_id": city}) // "country": countries[city[:len(city)-2]], "name": city.Name
	// }
	// var st []map[string]interface{}
	// for _, trip := range u.GetSuggestedTrips() {
	// 	st = append(st, map[string]interface{}{"trip_id": trip}) // "country": countries[city[:len(city)-2]], "name": city.Name
	// }
	var sc []string
	for _, city := range u.GetSuggestedCities() {
		sc = append(sc, city) // "country": countries[city[:len(city)-2]], "name": city.Name
	}
	var st []string
	for _, trip := range u.GetSuggestedTrips() {
		st = append(st, trip) // "country": countries[city[:len(city)-2]], "name": city.Name
	}
	var smt []string
	for _, trip := range u.GetSuggestedManualTrips() {
		smt = append(smt, trip) // "country": countries[city[:len(city)-2]], "name": city.Name
	}
	batch.
		Update(user, []firestore.Update{
			{Path: "needs_update", Value: false},
			{Path: "cat_params", Value: u.GetCats()},
			{Path: "mood_params", Value: u.GetMoods()},
			{Path: "pop", Value: u.GetPop()},
			{Path: "suggested_cities", Value: sc},
			{Path: "suggested_trips", Value: st},
			{Path: "suggested_manual_trips", Value: smt}})
	if len(yourTrips) > 0 {
		iter := client.Collection("trips").Where("user_id", "==", u.GetID()).Where("tag", "==", "yours").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			remove := false
			for _, trip := range yourTrips {
				if doc.Data()["city_id"].(string) == trip.CityID {
					remove = true
				}
			}
			if remove {
				batch.Delete(doc.Ref)
			}
		}
	}
	for _, trip := range yourTrips {
		// fmt.Println("Added your trip")
		var places []map[string]interface{}
		duration := 0.0
		for _, place := range trip.Places {
			places = append(places, map[string]interface{}{"coordinates": base.Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
			duration += place.Duration
		}
		tripRef := client.Collection("trips").Doc(trip.TripID)
		duration += trip.Time / 60.0
		tripFB := map[string]interface{}{
			"city_id":            trip.CityID, //include in addition to existing trips
			"user_id":            u.GetID(),
			"distance":           trip.Distance,
			"global_cat_params":  trip.Cat,
			"global_mood_params": trip.Mood,
			"pop":                trip.Pop,
			"sights":             places,
			"source":             "auto",
			"tag":                "yours",
			"time":               trip.Time,
			"walk_dist":          trip.WalkDist,
			"duration":           duration,
			"description":        "",
			"name":               "",
			"timestamp":          time.Now().Unix()}
		batch.Set(tripRef, tripFB)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	// fmt.Println("Updated")
}

func launchHTTPServer() {
	// add a function to handle sorting the places later
	http.HandleFunc("/make_your_trip", handleMakeYour)
	http.HandleFunc("/create_trip", handleCreate)
	http.HandleFunc("/edit_trip", handleEdit)
	http.HandleFunc("/remove_trip", handleRemove)
	http.HandleFunc("/get_match", handleMatch)
	http.HandleFunc("/get_name", handleName)
	http.HandleFunc("/get_images", handleImages)
	http.HandleFunc("/get_trip_images", handleTripImages)
	http.HandleFunc("/get_place_images", handlePlaceImages)
	// http.HandleFunc("/get_images_create", handleImageCreateTrip)
	fmt.Println("launching HTTP")
	http.ListenAndServe(":8080", nil)
}

// this handles http create trip request (cityID param), include a firebase read here
func handleMakeYour(w http.ResponseWriter, r *http.Request) {
	dict := r.URL.Query()
	fmt.Println(dict)
	cityID := ""
	if val, ok := dict["city_id"]; ok {
		cityID = val[0]
	}
	UID := ""
	if val, ok := dict["id"]; ok {
		UID = val[0]
	}
	tripID := xid.New().String()
	go makeYourAsync(cityID, tripID, UID)
	var _res = make(map[string]string)
	_res["trip_id"] = tripID
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
}

func makeYourAsync(cityID string, tripID string, UID string) {
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
	client, _ := firestore.NewClient(ctx, "***", opt)
	user := base.Users[UID]
	madeTrip := base.MakeYourTrip(cityID, tripID, user.GetCats(), user.GetMoods(), user.GetPop(), UID)
	tripRef := client.Collection("trips").Doc(tripID)
	var places []map[string]interface{}
	duration := 0.0
	for _, place := range madeTrip.Places {
		places = append(places, map[string]interface{}{"coordinates": base.Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
		duration += place.Duration
	}
	duration += madeTrip.Time / 60.0
	trip := map[string]interface{}{
		"city_id":            madeTrip.CityID, //include in addition to existing trips
		"user_id":            UID,
		"distance":           madeTrip.Distance,
		"global_cat_params":  madeTrip.Cat,
		"global_mood_params": madeTrip.Mood,
		"pop":                madeTrip.Pop,
		"sights":             places,
		"source":             "auto",
		"tag":                "yours",
		"time":               madeTrip.Time,
		"walk_dist":          madeTrip.WalkDist,
		"duration":           duration,
		"description":        "",
		"name":               "",
		"timestamp":          time.Now().Unix()}
	tripRef.Set(ctx, trip)
}

// this handles http adjust trip request (cityID, places param)
func handleCreate(w http.ResponseWriter, r *http.Request) {
	dict := r.URL.Query()
	fmt.Println(dict)
	cityID := ""
	if val, ok := dict["city_id"]; ok {
		cityID = val[0]
	}
	UID := ""
	if val, ok := dict["uid"]; ok {
		UID = val[0]
	}
	var placeIDs []string
	if val, ok := dict["place_ids[]"]; ok {
		// json.Unmarshal([]byte(val[0]), &placeIDs)
		placeIDs = val
	}
	name := ""
	if val, ok := dict["name"]; ok {
		name = val[0]
	}
	description := ""
	if val, ok := dict["description"]; ok {
		description = val[0]
	}
	privacy := ""
	if val, ok := dict["privacy"]; ok {
		privacy = val[0]
	}
	tripID := xid.New().String()
	go createAsync(cityID, tripID, UID, placeIDs, name, description, privacy)
	var _res = make(map[string]string)
	_res["trip_id"] = tripID
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
}

func createAsync(cityID string, tripID string, UID string, placeIDs []string, name string, description string, privacy string) {
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
	client, _ := firestore.NewClient(ctx, "***", opt)
	places := make([]base.Place, 0)
	for _, placeID := range placeIDs {
		places = append(places, base.Places[cityID][placeID])
	}
	createdTrip := base.CreateTrip(cityID, tripID, places, UID, privacy) // return trip id after writing to firestore
	tripRef := client.Collection("trips").Doc(tripID)
	var _places []map[string]interface{}
	duration := 0.0
	for _, place := range createdTrip.Places {
		_places = append(_places, map[string]interface{}{"coordinates": base.Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
		duration += place.Duration
	}
	duration += createdTrip.Time / 60.0
	trip := map[string]interface{}{
		"city_id":            createdTrip.CityID, //include in addition to existing trips
		"user_id":            UID,
		"distance":           createdTrip.Distance,
		"global_cat_params":  createdTrip.Cat,
		"global_mood_params": createdTrip.Mood,
		"pop":                createdTrip.Pop,
		"sights":             _places,
		"source":             "manual",
		"tag":                "created",
		"time":               createdTrip.Time,
		"walk_dist":          createdTrip.WalkDist,
		"duration":           duration,
		"description":        description,
		"name":               name,
		"privacy":            privacy,
		"likes":              0,
		"timestamp":          time.Now().Unix()}
	tripRef.Set(ctx, trip)
}

// this handles http adjust trip request (cityID, places param)
func handleEdit(w http.ResponseWriter, r *http.Request) {
	dict := r.URL.Query()
	fmt.Println(dict)
	cityID := ""
	if val, ok := dict["city_id"]; ok {
		cityID = val[0]
	}
	UID := ""
	if val, ok := dict["uid"]; ok {
		UID = val[0]
	}
	tripID := ""
	if val, ok := dict["trip_id"]; ok {
		tripID = val[0]
	}
	var placeIDs []string
	if val, ok := dict["place_ids[]"]; ok {
		placeIDs = val
	}
	name := ""
	if val, ok := dict["name"]; ok {
		name = val[0]
	}
	description := ""
	if val, ok := dict["description"]; ok {
		description = val[0]
	}
	privacy := ""
	if val, ok := dict["privacy"]; ok {
		privacy = val[0]
	}
	var _res = make(map[string]string)
	if base.Trips[cityID][tripID].Tag == "created" && base.Trips[cityID][tripID].UserID == UID {
		go editAsync(cityID, tripID, tripID, UID, placeIDs, name, description, privacy)
		_res["trip_id"] = tripID
	} else {
		newTripID := xid.New().String()
		go editAsync(cityID, tripID, newTripID, UID, placeIDs, name, description, privacy)
		_res["trip_id"] = newTripID
	}
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
}

func editAsync(cityID string, tripID string, newTripID string, UID string, placeIDs []string, name string, description string, privacy string) {
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
	client, _ := firestore.NewClient(ctx, "***", opt)
	places := make([]base.Place, 0)
	for _, placeID := range placeIDs {
		places = append(places, base.Places[cityID][placeID])
	}
	editedTrip := base.EditTrip(cityID, tripID, newTripID, places, UID, privacy) // return trip id after writing to firestore
	var _places []map[string]interface{}
	duration := 0.0
	for _, place := range editedTrip.Places {
		_places = append(_places, map[string]interface{}{"coordinates": base.Places[place.CityID][place.PlaceID].Coords, "name": place.Name, "place_id": place.PlaceID})
		duration += place.Duration
	}
	duration += editedTrip.Time / 60.0
	if tripID == newTripID {
		tripRef := client.Collection("trips").Doc(tripID)
		tripRef.Update(ctx, []firestore.Update{
			{
				Path:  "name",
				Value: name,
			},
			{
				Path:  "description",
				Value: description,
			},
			{
				Path:  "distance",
				Value: editedTrip.Distance,
			},
			{
				Path:  "global_cat_params",
				Value: editedTrip.Cat,
			},
			{
				Path:  "global_mood_params",
				Value: editedTrip.Mood,
			},
			{
				Path:  "pop",
				Value: editedTrip.Pop,
			},
			{
				Path:  "sights",
				Value: _places,
			},
			{
				Path:  "source",
				Value: "manual",
			},
			{
				Path:  "tag",
				Value: "created",
			},
			{
				Path:  "time",
				Value: editedTrip.Time,
			},
			{
				Path:  "walk_dist",
				Value: editedTrip.WalkDist,
			},
			{
				Path:  "duration",
				Value: duration,
			},
			{
				Path:  "privacy",
				Value: privacy,
			},
		})
	} else {
		trip := map[string]interface{}{
			"city_id":            cityID,
			"user_id":            UID,
			"name":               name,
			"description":        description,
			"distance":           editedTrip.Distance,
			"global_cat_params":  editedTrip.Cat,
			"global_mood_params": editedTrip.Mood,
			"pop":                editedTrip.Pop,
			"sights":             _places,
			"source":             "manual",
			"tag":                "created",
			"time":               editedTrip.Time,
			"walk_dist":          editedTrip.WalkDist,
			"duration":           duration,
			"privacy":            privacy,
			"likes":              0,
			"timestamp":          time.Now().Unix()}
		tripRef := client.Collection("trips").Doc(newTripID)
		tripRef.Set(ctx, trip)
	}
}

// this handles http adjust trip request (cityID, places param)
func handleRemove(w http.ResponseWriter, r *http.Request) { // is this safe
	dict := r.URL.Query()
	fmt.Println(dict)
	cityID := ""
	if val, ok := dict["city_id"]; ok {
		cityID = val[0]
	}
	tripID := ""
	if val, ok := dict["trip_id"]; ok {
		tripID = val[0]
	}
	removeAsync(tripID)
	delete(base.Trips[cityID], tripID)
	var _res = make(map[string]bool)
	_res["ok"] = true
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
}

func removeAsync(tripID string) {
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
		"client_x509_cert_url": ""
	  }`))
	ctx := context.Background()
	client, _ := firestore.NewClient(ctx, "***", opt)
	tripRef := client.Collection("trips").Doc(tripID)
	tripRef.Delete(ctx)
}

// this handles http get match request
func handleMatch(w http.ResponseWriter, r *http.Request) {
	dict := r.URL.Query()
	mc := 0.0
	cityID, uid := "", ""
	if val, ok := dict["city_id"]; ok {
		cityID = val[0]
	}
	if val, ok := dict["id"]; ok {
		uid = val[0]
	}
	if u, ok := base.Users[uid]; ok {
		if val, ok := dict["trip_id"]; ok {
			key := val[0]
			trip := base.Trips[cityID][key]
			mc = trip.GetTripMC(u.GetCats(), u.GetMoods(), u.GetPop())
			var _res = make(map[string]float64)
			_res["mc"] = mc
			res, _ := json.Marshal(_res)
			w.Write([]byte(res))
			return
		}
		if val, ok := dict["place_id"]; ok {
			key := val[0]
			place := base.Places[cityID][key]
			mc = place.GetPlaceMC(u.GetCats(), u.GetMoods(), u.GetPop())
			var _res = make(map[string]float64)
			_res["mc"] = mc
			res, _ := json.Marshal(_res)
			w.Write([]byte(res))
			return
		}
		city := base.Cities[cityID]
		mc = city.GetCityMC(u.GetCats(), u.GetMoods())
		var _res = make(map[string]float64)
		_res["mc"] = mc
		res, _ := json.Marshal(_res)
		w.Write([]byte(res))
		return
	}
	var _res = make(map[string]float64)
	_res["mc"] = 50
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
	return
}

// this handles http get user name request
func handleName(w http.ResponseWriter, r *http.Request) {
	dict := r.URL.Query()
	uid := ""
	if val, ok := dict["uid"]; ok {
		uid = val[0]
	}
	name := base.Users[uid].GetName()
	var _res = make(map[string]string)
	_res["name"] = name
	res, _ := json.Marshal(_res)
	w.Write([]byte(res))
	return
}

// Float64ToByte conversion
func Float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// this handles http get images request
func handleImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return
	}
	bkt := client.Bucket("pyxy_place_images")
	dict := r.URL.Query()
	key := ""
	if val, ok := dict["city_id"]; ok {
		key = val[0]
	}
	var urls = make(map[string]map[string]string)
	var urlsCity = make(map[string]string)
	for _, place := range base.Places[key] {
		query := &storage.Query{Prefix: key + "/" + place.PlaceID}
		it := bkt.Objects(ctx, query)
		item, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		urlsCity[place.PlaceID] = item.MediaLink
	}
	urls["imgs"] = urlsCity

	imgs, err := json.Marshal(urls)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(urls)
	w.Write([]byte(imgs))
}

// this handles http get images request
func handleTripImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return
	}
	bkt := client.Bucket("pyxy_place_images")
	dict := r.URL.Query()
	key := ""
	if val, ok := dict["city_id"]; ok {
		key = val[0]
	}
	trip := ""
	if val, ok := dict["trip_id"]; ok {
		trip = val[0]
	}
	var urls = make(map[string]map[string]string)
	var urlsTrip = make(map[string]string)
	for _, place := range base.Trips[key][trip].Places {
		query := &storage.Query{Prefix: key + "/" + place.PlaceID}
		it := bkt.Objects(ctx, query)
		item, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		urlsTrip[place.PlaceID] = item.MediaLink
	}
	urls["imgs"] = urlsTrip
	imgs, err := json.Marshal(urls)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(imgs))
}

// this handles http get images request
func handlePlaceImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return
	}
	bkt := client.Bucket("pyxy_place_images")
	dict := r.URL.Query()
	key := ""
	if val, ok := dict["city_id"]; ok {
		key = val[0]
	}
	placeID := ""
	if val, ok := dict["place_id"]; ok {
		placeID = val[0]
	}
	var urls = make(map[string][]string)
	var urlsPlace = make([]string, 0)
	query := &storage.Query{Prefix: key + "/" + placeID}
	it := bkt.Objects(ctx, query)
	i := 0
	for {
		item, err := it.Next()
		if err == iterator.Done || i == 5 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		urlsPlace = append(urlsPlace, item.MediaLink)
		i++
	}
	urls["imgs"] = urlsPlace
	imgs, err := json.Marshal(urls)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(urls)
	w.Write([]byte(imgs))
}
