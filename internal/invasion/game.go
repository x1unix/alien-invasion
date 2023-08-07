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
	aliens []*Alien
	cities mapfile.Nodes

	tickLimit       uint
	tickCount       uint
	isAliensChanged bool
	intersections   map[string]alienRef
}

func NewGame(aliensCount int, tickLimit uint, cities mapfile.Nodes) *Game {
	aliens := GenerateAliens(aliensCount, cities)

	return &Game{
		aliens:    aliens,
		cities:    cities,
		tickLimit: tickLimit,
	}
}

// Cities returns remaining alive cities.
func (g *Game) Cities() mapfile.Nodes {
	return g.cities
}

// Tick performs one game move and returns true if game is finished.
func (g *Game) Tick() bool {
	if g.tickCount == 0 {
		g.beginGame()
	}

	if g.tickCount >= g.tickLimit {
		zap.L().Info("Game is finished, step limit reached")
		return false
	}

	g.walk()
	g.tickCount++
	return len(g.aliens) > 0
}

func (g *Game) beginGame() {
	g.intersections = make(map[string]alienRef)
	for i, alien := range g.aliens {
		zap.L().Info("alien landed",
			zap.Int("alien_id", alien.ID), zap.Stringer("city", alien.CurrentCity),
		)

		g.intersections[alien.CurrentCity.Name] = alienRef{
			index: i,
			alien: alien,
		}
	}
}

func (g *Game) walk() {
	for i, alien := range g.aliens {
		if alien == nil {
			// skip if alien was removed during the turn.
			continue
		}

		if alien.IsStuck() {
			zap.L().Info("alien stuck and can't move",
				zap.Int("alien_id", alien.ID), zap.Stringer("city", alien.CurrentCity),
			)

			continue
		}

		originCity := alien.CurrentCity.Name
		alien.MoveNext()

		// Remove previous place where alien was before move.
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

		zap.L().Info("aliens killed each other",
			zap.Int("alien1", alien.ID),
			zap.Int("alien2", prevAlien.alien.ID),
			zap.String("city", newCity),
		)

		fmt.Printf(
			"%s has been destroyed by alien %d and alien %d!\n",
			newCity, alien.ID, prevAlien.alien.ID,
		)

		// Remove related city
		g.removeCity(alien.CurrentCity)
		delete(g.intersections, newCity)

		// Mark aliens to be removed
		g.aliens[prevAlien.index] = nil
		g.aliens[i] = nil
		g.isAliensChanged = true
	}

	g.cleanup()
	zap.L().Info("end of step",
		zap.Int("aliens_alive", len(g.aliens)),
		zap.Uint("step_count", g.tickCount))
}

func (g *Game) cleanup() {
	// reset per-step intersections state
	g.intersections = make(map[string]alienRef)

	if !g.isAliensChanged {
		return
	}

	// cleanup empty aliens slots
	newAliens := make([]*Alien, 0, len(g.aliens))
	for _, alien := range g.aliens {
		if alien == nil {
			continue
		}

		newAliens = append(newAliens, alien)
	}

	g.aliens = newAliens
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

	delete(g.cities, city.Name)
}
