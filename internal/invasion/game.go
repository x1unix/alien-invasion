package invasion

import (
	"log"
	"math/rand"

	"github.com/x1unix/alien-invasion/internal/mapfile"
)

type Game struct {
	Aliens []*Alien
	Cities mapfile.Cities
}

func NewGame(aliensCount int, cities mapfile.Cities) *Game {
	aliens := GenerateAliens(aliensCount, cities)

	return &Game{Aliens: aliens, Cities: cities}
}

func (g *Game) Walk() {
	for _, alien := range g.Aliens {
		nextDirection := getAlienDirection(alien)
		if nextDirection == nil {
			// TODO: remove from pool
			log.Printf(
				"Alien#%d has nowhere to go and stuck at %s",
				alien.ID, alien.CurrentCity,
			)
			continue
		}

		alien.CurrentCity = nextDirection
		log.Printf("Alien#%d went to %s", alien.ID, alien.CurrentCity)
	}
}

func getAlienDirection(a *Alien) *mapfile.Node {
	nextCities := a.CurrentCity.Directions()
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
