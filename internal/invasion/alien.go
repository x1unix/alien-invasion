package invasion

import (
	"math/rand"

	"github.com/x1unix/alien-invasion/internal/mapfile"
)

type Alien struct {
	ID          int
	Alive       bool
	CurrentCity *mapfile.Node
}

func GenerateAliens(count int, cities mapfile.Cities) []*Alien {
	aliens := make([]*Alien, count)
	landingCities := getRandomInvasionCities(count, cities)

	for i, city := range landingCities {
		aliens[i] = &Alien{
			ID:          i,
			CurrentCity: city,
		}
	}

	return aliens
}

func getRandomInvasionCities(n int, cities mapfile.Cities) []*mapfile.Node {
	elems := cities.AsSlice()
	result := make([]*mapfile.Node, n)

	for i := 0; i < n; i++ {
		j := rand.Intn(len(elems) - 1)
		result[i] = elems[j]
		elems = append(elems[:j], elems[j+1:]...)
	}

	return result
}
