package base

type Population struct {
	Routes []Route
}

func (a *Population) InitEmpty(pSize int) {
	a.Routes = make([]Route, pSize)
}

func (a *Population) InitPopulation(pSize int, tm TripManager) {
	a.Routes = make([]Route, pSize)
	for i := 0; i < pSize; i++ {
		nT := Route{}
		nT.InitRoutePlaces(tm)
		a.SaveRoute(i, nT)
	}
}

func (a *Population) SaveRoute(i int, t Route) {
	a.Routes[i] = t
}

func (a *Population) GetRoute(i int) *Route {
	return &a.Routes[i]
}

func (a *Population) PopulationSize() int {
	return len(a.Routes)
}

func (a *Population) GetFittest() *Route {
	fittest := a.Routes[0]
	// Loop through all Trips taken by population and determine the fittest
	for i := 0; i < a.PopulationSize(); i++ {
		// log.Println("Current Trip: ", i)
		// fmt.Println(a.GetRoute(i).Fitness())
		if fittest.Fitness() <= a.GetRoute(i).Fitness() {
			fittest = *a.GetRoute(i)
		}
	}
	// fmt.Printf("best = %v", fittest.Fitness())
	// fmt.Println(fittest.GetPlaces())
	return &fittest
}
