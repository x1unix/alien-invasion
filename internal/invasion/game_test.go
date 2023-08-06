package invasion

import (
	"github.com/stretchr/testify/require"
	"github.com/x1unix/alien-invasion/internal/mapfile"
	"testing"
)

func TestGame_removeCity(t *testing.T) {
	m, err := mapfile.ReadFile("../../map.txt")
	require.NoError(t, err)

	game := NewGame(1, 10, m.Cities)
	//game.Tick()

	c, ok := game.Cities["Whiterun"]
	game.Aliens[0].CurrentCity = game.Cities["Windhelm"]
	require.True(t, ok)
	game.removeCity(c)
	t.Log()
}
