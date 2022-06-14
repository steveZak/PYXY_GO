package base

import (
	"math"
	"math/rand"
	"sort"
)

// City : city struct
type City struct {
	cityID string
	mc     float64
	cat    []float64
	mood   []float64
}

func GenerateCity(cityID string, cat []interface{}, mood []interface{}) City {
	c := City{}
	c.cityID = cityID
	for _, v := range cat {
		c.cat = append(c.cat, v.(float64))
	}
	for _, v := range mood {
		c.mood = append(c.mood, v.(float64))
	}
	return c
}

func (c *City) GetCityID() string {
	return c.cityID
}

// cosine similarity
func (c *City) GetCityMC(cat []float64, mood []float64) float64 {
	sumUC, sumUU, sumCC := 0.0, 0.0, 0.0
	for i := 0; i < len(cat); i++ {
		sumUC += cat[i] * c.cat[i]
		sumUU += cat[i] * cat[i]
		sumCC += c.cat[i] * c.cat[i]
	}
	kC := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	sumUC, sumUU, sumCC = 0.0, 0.0, 0.0
	for i := 0; i < len(mood); i++ {
		sumUC += mood[i] * c.mood[i]
		sumUU += mood[i] * mood[i]
		sumCC += c.mood[i] * c.mood[i]
	}
	kM := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	return 0.8*(kC+1.0)/2.0 + 0.2*(kM+1.0)/2.0
}

func RankCities(cat []float64, mood []float64) []City { //consider: match countries
	var c []City
	for _, city := range Cities { // untested
		c = append(c, city)
	}
	c[0].mc = c[0].GetCityMC(cat, mood)
	cSort := []City{}
	for i := 1; i < len(c); i++ {
		c[i].mc = c[i].GetCityMC(cat, mood)
		j := sort.Search(len(cSort), func(j int) bool { return cSort[j].mc < c[i].mc })
		if j == len(cSort) {
			cSort = append(cSort, []City{c[i]}...)
			continue
		}
		if j == 0 {
			cSort = append([]City{c[i]}, cSort...)
			continue
		}
		cSort = append(cSort[0:j], append([]City{c[i]}, cSort[j:]...)...)
	}
	return cSort
}

// ShuffleCities : if params are default
func ShuffleCities(in []City) []City {
	out := make([]City, len(in), cap(in))
	perm := rand.Perm(len(in))
	for i, v := range perm {
		out[v] = in[i]
	}
	return out
}
