package base

// TripManager : Contains list of of cities to be visited
type TripManager struct {
	destPlaces []Place
}

// NewTripManager : Initialize TripManager
func (a *TripManager) NewTripManager(routeSize int) {
	a.destPlaces = make([]Place, 0, routeSize)
}

func (a *TripManager) AddPlace(p Place) {
	a.destPlaces = append(a.destPlaces, p)
}

func (a *TripManager) GetPlace(i int) Place {
	return a.destPlaces[i]
}

func (a *TripManager) NumberOfPlaces() int {
	return len(a.destPlaces)
}
