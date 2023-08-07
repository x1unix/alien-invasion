package invasion

import (
	"fmt"

	"github.com/x1unix/alien-invasion/internal/mapfile"
	"go.uber.org/zap"
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
		zap.S().Infof("Game is finished, step limit reached")
		return false
	}

	g.walk()
	g.tickCount++
	return len(g.Aliens) > 0
}

func (g *Game) beginGame() {
	g.intersections = make(map[string]alienRef)
	for i, alien := range g.Aliens {
		zap.S().Infof("Alien#%d landed on %s", alien.ID, alien.CurrentCity)
		g.intersections[alien.CurrentCity.Name] = alienRef{
			index: i,
			alien: alien,
		}
	}
}

func (g *Game) walk() {
	for i, alien := range g.Aliens {
		if alien == nil {
			continue
		}
		if alien.IsStuck() {
			zap.S().Infof("Alien#%d stuck in %s", alien.ID, alien.CurrentCity.String())
			continue
		}

		originCity := alien.CurrentCity.Name
		alien.MoveNext()

		// Remove previous place where alien was.
		delete(g.intersections, originCity)

		newCity := alien.CurrentCity.Name
		prevAlien, ok := g.intersections[newCity]
		if !ok {
			// Remember place where alien arrived.
			g.intersections[newCity] = alienRef{
				index: i,
				alien: alien,
			}
			continue
		}

		zap.S().Infof(
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
	zap.L().Info("end of step",
		zap.Int("aliens_alive", len(g.Aliens)),
		zap.Uint("step_count", g.tickCount))
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
	zap.L().Debug("remove city",
		zap.Stringer("city", city),
		zap.Stringers("siblings", []fmt.Stringer{city.South, city.North, city.West, city.East}))

	if city.South != nil {
		city.South.North = nil
		city.South = nil
	}
	if city.North != nil {
		city.North.South = nil
		city.North = nil
	}
	if city.West != nil {
		city.West.East = nil
		city.West = nil
	}
	if city.East != nil {
		city.East.West = nil
		city.East = nil
	}

	delete(g.Cities, city.Name)
}
