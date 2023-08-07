package invasion

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/x1unix/alien-invasion/internal/invasion/mock"
	"github.com/x1unix/alien-invasion/internal/mapfile"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestGame_Tick(t *testing.T) {
	cases := map[string]struct {
		aliensCount int
		tickLimit   uint
		mapFile     string

		prepareRand    func(ctrl *mock.MockRand)
		prepareGame    func(game *Game)
		checkGameState func(t *testing.T, game *Game)
	}{
		"respects tick limit": {
			tickLimit:   2,
			aliensCount: 1,
			mapFile:     "skyrim.txt",
			checkGameState: func(t *testing.T, game *Game) {
				require.LessOrEqual(t, uint(2), game.tickCount)
			},
		},
		"should destroy city on intersect": {
			aliensCount: 2,
			tickLimit:   10000,
			mapFile:     "collapse.txt",
			prepareGame: func(game *Game) {
				game.aliens = []*Alien{
					NewAlien(0, game.cities["left"]),
					NewAlien(1, game.cities["right"]),
				}
			},
			checkGameState: func(t *testing.T, game *Game) {
				wantCities := []string{"left", "right"}
				got := game.cities.Names()
				sort.Strings(got)
				require.Equal(t, wantCities, got)
				require.Empty(t, game.aliens)
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			f, err := mapfile.ReadFile(filepath.Join("testdata", c.mapFile))
			require.NoError(t, err, "failed to read test map file")
			game := NewGame(c.aliensCount, c.tickLimit, f.Cities)

			zap.ReplaceGlobals(zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel)))
			t.Cleanup(zap.ReplaceGlobals(zap.NewNop()))

			if c.prepareRand != nil {
				t.Cleanup(func() {
					rnd = originalRand
				})

				ctrl := gomock.NewController(t)
				r := mock.NewMockRand(ctrl)
				c.prepareRand(r)
				rnd = r
			}

			if c.prepareGame != nil {
				c.prepareGame(game)
			}

			for game.Tick() {
			}

			if c.checkGameState != nil {
				c.checkGameState(t, game)
			}
		})
	}
}
