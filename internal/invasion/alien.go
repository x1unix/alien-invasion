package invasion

import (
	"github.com/x1unix/alien-invasion/internal/mapfile"
	"go.uber.org/zap"
)

type Alien struct {
	ID          int
	MoveCount   int
	CurrentCity *mapfile.Node

	stuck bool
}

func NewAlien(id int, city *mapfile.Node) *Alien {
	return &Alien{
		ID:          id,
		CurrentCity: city,
	}
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
	randIndex := rnd.Intn(len(nextCities) - 1)
	return nextCities[randIndex]
}

// GenerateAliens generates a list of aliens that will land on random city.
func GenerateAliens(count int, cities []*mapfile.Node) []*Alien {
	aliens := make([]*Alien, count)
	landingCities := getRandomInvasionCities(count, cities)

	for i, city := range landingCities {
		aliens[i] = NewAlien(i, city)
	}

	return aliens
}

func getRandomInvasionCities(n int, cities []*mapfile.Node) []*mapfile.Node {
	elems := append([]*mapfile.Node{}, cities...)
	result := make([]*mapfile.Node, n)

	for i := 0; i < n; i++ {
		j := len(elems) - 1
		if j > 0 {
			j = rnd.Intn(len(elems) - 1)
		}
		result[i] = elems[j]
		elems = append(elems[:j], elems[j+1:]...)
	}

	return result
}
