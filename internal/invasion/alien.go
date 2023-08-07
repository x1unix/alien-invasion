package invasion

import (
	"math/rand"

	"github.com/x1unix/alien-invasion/internal/mapfile"
	"go.uber.org/zap"
)

type Alien struct {
	ID          int
	MoveCount   int
	CurrentCity *mapfile.Node

	stuck bool
}

func (a *Alien) IsStuck() bool {
	return a.stuck
}

func (a *Alien) MoveNext() {
	nextDirection := getAlienDirection(a)
	if nextDirection == nil {
		a.stuck = true
		zap.L().Info("alien stuck",
			zap.Int("alien_id", a.ID),
			zap.Stringer("city", a.CurrentCity),
		)
		return
	}

	a.MoveTo(nextDirection)
}

func (a *Alien) MoveTo(city *mapfile.Node) {
	a.MoveCount++
	a.CurrentCity = city
	zap.L().Info("alien moved",
		zap.Int("alien_id", a.ID),
		zap.Stringer("city", a.CurrentCity),
	)
}

func getAlienDirection(a *Alien) *mapfile.Node {
	nextCities := a.CurrentCity.Directions()
	zap.L().Debug("checking alien available directions",
		zap.Int("alien_id", a.ID),
		zap.Stringers("directions", nextCities))

	switch len(nextCities) {
	case 0:
		return nil
	case 1:
		return nextCities[0]
	}

	// Return random next city
	randIndex := rand.Intn(len(nextCities) - 1)
	return nextCities[randIndex]
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
