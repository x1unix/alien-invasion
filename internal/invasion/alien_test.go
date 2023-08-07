package invasion

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/x1unix/alien-invasion/internal/invasion/mock"
	"github.com/x1unix/alien-invasion/internal/mapfile"
)

var originalRand = rnd

//go:generate mockgen -package=mock -destination=mock/rand.go github.com/x1unix/alien-invasion/internal/invasion Rand
func TestAlien_MoveNext(t *testing.T) {
	cases := map[string]struct {
		alien       *Alien
		shouldStuck bool
		expectCity  *mapfile.Node

		prepareRand func(r *mock.MockRand)
	}{
		"should become stuck if no directions available": {
			shouldStuck: true,
			expectCity:  &mapfile.Node{Name: "city1"},
			alien: NewAlien(0, &mapfile.Node{
				Name: "city1",
			}),
		},
		"should pick remaining option": {
			expectCity: &mapfile.Node{Name: "city2"},
			alien: NewAlien(0, &mapfile.Node{
				Name: "city1",
				North: &mapfile.Node{
					Name: "city2",
				},
			}),
		},
		"should pick random city": {
			expectCity: &mapfile.Node{Name: "city3"},
			alien: NewAlien(0, &mapfile.Node{
				Name: "city1",
				East: &mapfile.Node{
					Name: "city2",
				},
				West: &mapfile.Node{
					Name: "city3",
				},
				South: &mapfile.Node{
					Name: "city4",
				},
				North: &mapfile.Node{
					Name: "city5",
				},
			}),
			prepareRand: func(r *mock.MockRand) {
				r.EXPECT().Intn(3).Return(1)
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			t.Cleanup(func() {
				rnd = originalRand
			})

			ctrl := gomock.NewController(t)
			r := mock.NewMockRand(ctrl)
			rnd = r

			if c.prepareRand != nil {
				c.prepareRand(r)
			}

			c.alien.MoveNext()
			require.Equal(t, c.shouldStuck, c.alien.IsStuck())
			require.Equal(t, c.expectCity, c.alien.CurrentCity)
		})
	}
}

func TestGenerateAliens(t *testing.T) {
	cities := []*mapfile.Node{
		{Name: "city1"},
		{Name: "city2"},
		{Name: "city3"},
		{Name: "city4"},
	}

	expect := []*Alien{
		NewAlien(0, &mapfile.Node{Name: "city1"}),
		NewAlien(1, &mapfile.Node{Name: "city3"}),
	}

	t.Cleanup(func() {
		rnd = originalRand
	})

	ctrl := gomock.NewController(t)
	r := mock.NewMockRand(ctrl)
	rnd = r

	r.EXPECT().Intn(3).Return(0)
	r.EXPECT().Intn(2).Return(1)
	got := GenerateAliens(2, cities)

	require.Equal(t, expect, got)
}

func TestAlien_IsStuck(t *testing.T) {
	a := Alien{stuck: true}
	require.Equal(t, a.stuck, a.IsStuck())
}

func TestAlien_MoveTo(t *testing.T) {
	newCity := &mapfile.Node{
		Name: "foo",
	}

	alien := NewAlien(1, nil)
	alien.MoveTo(newCity)

	require.Equal(t, 1, alien.MoveCount)
	require.Equal(t, newCity, alien.CurrentCity)
}
