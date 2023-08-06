package invasion

import (
	"log"

	"github.com/x1unix/alien-invasion/internal/mapfile"
)

type alienRef struct {
	index int
	alien *Alien
}

type Game struct {
	Aliens []*Alien
	Cities mapfile.Cities

	tickLimit       uint
	tickCount       uint
	isAliensChanged bool
	isFinished      bool
	intersections   map[string]alienRef
}

func NewGame(aliensCount int, tickLimit uint, cities mapfile.Cities) *Game {
	aliens := GenerateAliens(aliensCount, cities)

	return &Game{
		Aliens:    aliens,
		Cities:    cities,
		tickLimit: tickLimit,
	}
}

// Tick performs one game move and returns true if game is finished
func (g *Game) Tick() bool {
	if g.tickCount == 0 {
		g.beginGame()
	}

	if g.tickCount >= g.tickLimit {
		return false
	}

	g.walk()
	g.tickCount++
	return len(g.Aliens) > 0
}

func (g *Game) beginGame() {
	g.intersections = make(map[string]alienRef)
	for i, alien := range g.Aliens {
		log.Printf("Alien#%d landed on %s", alien.ID, alien.CurrentCity)
		g.intersections[alien.CurrentCity.Name] = alienRef{
			index: i,
			alien: alien,
		}
	}
}

func (g *Game) walk() {
	for i, alien := range g.Aliens {
		if alien == nil || alien.IsStuck() {
			continue
		}

		alien.MoveNext()
		newCity := alien.CurrentCity.Name
		prevAlien, ok := g.intersections[newCity]
		if !ok {
			g.intersections[newCity] = alienRef{
				index: i,
				alien: alien,
			}
			continue
		}

		log.Printf(
			"Alien#%d and Alien#%d killed each other in %s",
			alien.ID, prevAlien.alien.ID, newCity,
		)

		// Remove related city
		g.removeCity(alien.CurrentCity)
		delete(g.intersections, newCity)

		// Mark aliens to be removed
		g.Aliens[prevAlien.index] = nil
		g.Aliens[i] = nil
		g.isAliensChanged = true
	}

	g.cleanup()
	log.Println("End of step")
}

func (g *Game) cleanup() {
	// reset per-step intersections state
	g.intersections = make(map[string]alienRef)

	if !g.isAliensChanged {
		return
	}

	// cleanup empty aliens slots
	newAliens := make([]*Alien, 0, len(g.Aliens))
	for _, alien := range g.Aliens {
		if alien == nil {
			continue
		}

		newAliens = append(newAliens, alien)
	}

	g.Aliens = newAliens
	g.isAliensChanged = false
}

func (g *Game) removeCity(city *mapfile.Node) {
	// Assume that each sibling is correctly connected.
	log.Println("DEBUG: removeCity ", city.String())
	log.Println("DEBUG: siblings ",
		city.South.String(), city.North.String(), city.West.String(), city.East.String())

	if city.South != nil {
		city.South.North = nil
		city.South = nil
	}
	if city.North != nil {
		city.North.South = nil
		city.North = nil
	}
	if city.West != nil {
		log.Printf("Debug: removeCity(%s) West(%s).East(%s)", city, city.West, city.East)
		city.West.East = nil
		city.West = nil
	}
	if city.East != nil {
		log.Printf("Debug: removeCity(%s) East(%s).West(%s)", city, city.East, city.West)
		city.East.West = nil
		city.East = nil
	}

	log.Println("DEBUG: remove from cities", city.String())
	delete(g.Cities, city.Name)
}
