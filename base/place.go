package base

import (
	"math"
	"math/rand"
	"sort"
)

// Place : place struct
type Place struct {
	CityID   string             `json:"cityID"`
	PlaceID  string             `json:"placeID"`
	Order    int                `json:"order"`
	Name     string             `json:"name"`
	MC       float64            `json:"mc"`
	Cat      []float64          `json:"cat"`
	Mood     []float64          `json:"mood"`
	Pop      float64            `json:"pop"`
	Coords   map[string]float64 `json:"coordinates"`
	Duration float64            `json:"duration"`
}

func GeneratePlace(CityID string, PlaceID string, Order int, Name string, Cat []interface{}, Mood []interface{}, Pop float64, Coords map[string]interface{}, Duration float64) Place {
	t := Place{}
	t.CityID = CityID
	t.PlaceID = PlaceID
	t.Order = Order
	t.Name = Name
	for _, v := range Cat {
		t.Cat = append(t.Cat, v.(float64))
	}
	for _, v := range Mood {
		t.Mood = append(t.Mood, v.(float64))
	}
	t.Pop = Pop
	var coords = make(map[string]float64)
	for k, v := range Coords {
		coords[k] = v.(float64)
	}
	t.Coords = coords
	t.Duration = Duration
	return t
}

func ImportPlace(CityID string, PlaceID string, Name string, Cat []interface{}, Mood []interface{}, Pop float64) Place {
	t := Place{}
	t.CityID = CityID
	t.PlaceID = PlaceID
	t.Name = Name
	for _, v := range Cat {
		t.Cat = append(t.Cat, v.(float64))
	}
	for _, v := range Mood {
		t.Mood = append(t.Mood, v.(float64))
	}
	t.Pop = Pop
	return t
}

// GetPlaceMC cosine similarity
func (p *Place) GetPlaceMC(Cat []float64, Mood []float64, Pop float64) float64 {
	sumUC, sumUU, sumCC := 0.0, 0.0, 0.0
	for i := 0; i < len(Cat); i++ {
		sumUC += Cat[i] * p.Cat[i]
		sumUU += Cat[i] * Cat[i]
		sumCC += p.Cat[i] * p.Cat[i]
	}
	kC := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	sumUC, sumUU, sumCC = 0.0, 0.0, 0.0
	for i := 0; i < len(Mood); i++ {
		sumUC += Mood[i] * p.Mood[i]
		sumUU += Mood[i] * Mood[i]
		sumCC += p.Mood[i] * p.Mood[i]
	}
	kM := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	kP := 1 - math.Abs(Pop-p.Pop)
	return 0.75*(kC+1.0)/2.0 + 0.15*(kM+1.0)/2.0 + 0.1*kP // (kP+1.0)/2.0
}

func RankPlaces(p []Place, Cat []float64, Mood []float64, Pop float64, limit int) []Place { //consider: match countries
	p[0].MC = p[0].GetPlaceMC(Cat, Mood, Pop)
	pSort := []Place{p[0]}
	for i := 1; i < len(p); i++ {
		p[i].MC = p[i].GetPlaceMC(Cat, Mood, Pop)
		j := sort.Search(len(pSort), func(j int) bool { return pSort[j].MC < p[i].MC })
		if j == len(pSort) {
			pSort = append(pSort, []Place{p[i]}...)
			continue
		}
		if j == 0 {
			pSort = append([]Place{p[i]}, pSort...)
			continue
		}
		pSort = append(pSort[0:j], append([]Place{p[i]}, pSort[j:]...)...)
	}
	return pSort[:limit]
}

func (p *Place) getPlaceMCCat(Cat []float64) float64 {
	sumUC, sumUU, sumCC := 0.0, 0.0, 0.0
	for i := 0; i < len(Cat); i++ {
		sumUC += Cat[i] * p.Cat[i]
		sumUU += Cat[i] * Cat[i]
		sumCC += p.Cat[i] * p.Cat[i]
	}
	kC := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	return (kC + 1.0) / 2.0
}

func (p *Place) getPlaceMCMood(Mood []float64) float64 {
	sumUC, sumUU, sumCC := 0.0, 0.0, 0.0
	for i := 0; i < len(Mood); i++ {
		sumUC += Mood[i] * p.Mood[i]
		sumUU += Mood[i] * Mood[i]
		sumCC += p.Mood[i] * p.Mood[i]
	}
	kM := sumUC / (math.Sqrt(sumUU) * math.Sqrt(sumCC))
	return (kM + 1.0) / 2.0
}

func (p *Place) getPlaceMCPop(Pop float64) float64 {
	kP := 1 - math.Abs(Pop-p.Pop)
	return (kP + 1.0) / 2.0
}

func RankPlacesCat(p []Place, Cat []float64, limit int) []Place { //consider: match countries
	p[0].MC = p[0].getPlaceMCCat(Cat)
	pSort := []Place{p[0]}
	for i := 1; i < len(p); i++ {
		p[i].MC = p[i].getPlaceMCCat(Cat)
		j := sort.Search(len(pSort), func(j int) bool { return pSort[j].MC < p[i].MC })
		if j == len(pSort) {
			pSort = append(pSort, []Place{p[i]}...)
			continue
		}
		if j == 0 {
			pSort = append([]Place{p[i]}, pSort...)
			continue
		}
		pSort = append(pSort[0:j], append([]Place{p[i]}, pSort[j:]...)...)
	}
	return pSort[:limit]
}

func RankPlacesMood(p []Place, Mood []float64, limit int) []Place { //consider: match countries
	p[0].MC = p[0].getPlaceMCMood(Mood)
	pSort := []Place{p[0]}
	for i := 1; i < len(p); i++ {
		p[i].MC = p[i].getPlaceMCMood(Mood)
		j := sort.Search(len(pSort), func(j int) bool { return pSort[j].MC < p[i].MC })
		if j == len(pSort) {
			pSort = append(pSort, []Place{p[i]}...)
			continue
		}
		if j == 0 {
			pSort = append([]Place{p[i]}, pSort...)
			continue
		}
		pSort = append(pSort[0:j], append([]Place{p[i]}, pSort[j:]...)...)
	}
	return pSort[:limit]
}

func RankPlacesPop(p []Place, Pop float64, limit int) []Place { //consider: match countries
	p[0].MC = p[0].getPlaceMCPop(Pop)
	pSort := []Place{p[0]}
	for i := 1; i < len(p); i++ {
		p[i].MC = p[i].getPlaceMCPop(Pop)
		j := sort.Search(len(pSort), func(j int) bool { return pSort[j].MC < p[i].MC })
		if j == len(pSort) {
			pSort = append(pSort, []Place{p[i]}...)
			continue
		}
		if j == 0 {
			pSort = append([]Place{p[i]}, pSort...)
			continue
		}
		pSort = append(pSort[0:j], append([]Place{p[i]}, pSort[j:]...)...)
	}
	return pSort[:limit]
}

func ShufflePlaces(in []Place) []Place {
	out := make([]Place, len(in), cap(in))
	perm := rand.Perm(len(in))
	for i, v := range perm {
		out[v] = in[i]
	}
	return out
}
