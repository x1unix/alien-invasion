package invasion

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/x1unix/alien-invasion/internal/mapfile"
)

func TestGame_removeCity(t *testing.T) {
	m, err := mapfile.ReadFile("../../map.txt")
	require.NoError(t, err)

	game := NewGame(1, 10, m.Cities)
	//game.Tick()

	c, ok := game.cities["Whiterun"]
	game.aliens[0].CurrentCity = game.cities["Windhelm"]
	require.True(t, ok)
	game.removeCity(c)
	t.Log()
}

func TestGame(t *testing.T) {
	m, err := mapfile.ReadFile("../../cites.txt")
	require.NoError(t, err)

	game := NewGame(10, 10000, m.Cities)
	places := []string{
		"Red_Bank", "Elizabeth", "Newark",
		"Yonkers", "New_York", "Woodbridge",
		"Stamford", "Port_Chester", "Rye", "New_Rochelle",
	}

	for i, place := range places {
		game.aliens[i].CurrentCity = m.Cities[place]
	}

}
