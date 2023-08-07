package mapfile

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	cases := map[string]struct {
		src       string
		want      Structure
		wantError string
	}{
		"empty file": {
			src: "empty.txt",
		},
		"should contain directions": {
			src:       "wo-direction.txt",
			wantError: `missing directions for city "Foo"`,
		},
		"direction should not be empty": {
			src:       "empty-dir.txt",
			wantError: `invalid direction attribute: "north="`,
		},
		"direction attribute should contain only one value": {
			src:       "invalid-dir.txt",
			wantError: `invalid direction attribute: "north=foo=bar"`,
		},
		"should validate direction type": {
			src:       "unknown-dir.txt",
			wantError: `invalid city direction: "foo=bar"`,
		},
		"should correctly parse file": {
			src: "correct.txt",
			want: Structure{
				Cities: Cities{
					"Whiterun": {
						Name:  "Whiterun",
						East:  &Node{Name: "Windhelm"},
						West:  &Node{Name: "Markarth"},
						South: &Node{Name: "Falkreath"},
						North: &Node{Name: "Dawnstar"},
					},
					"Riften": {
						Name:  "Riften",
						West:  &Node{Name: "Falkreath"},
						North: &Node{Name: "Windhelm"},
					},
					"Windhelm": {
						Name:  "Windhelm",
						West:  &Node{Name: "Whiterun"},
						South: &Node{Name: "Riften"},
					},
					"Falkreath": {
						Name:  "Falkreath",
						North: &Node{Name: "Whiterun"},
						East:  &Node{Name: "Riften"},
					},
					"Markarth": {
						Name: "Markarth",
						East: &Node{Name: "Whiterun"},
					},
					"Dawnstar": {
						Name:  "Dawnstar",
						South: &Node{Name: "Whiterun"},
					},
				},
			},
		},
		"should return file read error": {
			src:       "not-exists.txt",
			wantError: "no such file or directory",
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			got, err := ReadFile(filepath.Join("testdata", c.src))
			if c.wantError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), c.wantError)
				return
			}

			require.NoError(t, err)
			compareCities(t, c.want.Cities, got.Cities)
		})
	}
}

// compareCities used to compare two cities map.
//
// Used instead of require.Equal due to recursive links.
func compareCities(t *testing.T, want, got Cities) {
	wantNames, gotNames := want.Names(), got.Names()
	sort.Strings(wantNames)
	sort.Strings(gotNames)
	require.Equal(t, wantNames, gotNames, "cities lists aren't identical")

	for name, wantCity := range want {
		gotCity := got[name]
		require.Equalf(
			t, wantCity.Name, gotCity.Name,
			"city name mismatch (key: %s)", name,
		)

		if wantCity.West != nil {
			require.Equalf(
				t, wantCity.West.Name, gotCity.West.Name,
				"city %q west direction mismatch", name,
			)
		} else {
			require.Nil(t, gotCity.West, "city %q west direction not nil", name)
		}

		if wantCity.East != nil {
			require.Equalf(
				t, wantCity.East.Name, gotCity.East.Name,
				"city %q east direction mismatch", name,
			)
		} else {
			require.Nil(t, gotCity.East, "city %q east direction not nil", name)
		}

		if wantCity.North != nil {
			require.Equalf(
				t, wantCity.North.Name, gotCity.North.Name,
				"city %q north direction mismatch", name,
			)
		} else {
			require.Nil(t, gotCity.North, "city %q north direction not nil", name)
		}

		if wantCity.South != nil {
			require.Equalf(
				t, wantCity.South.Name, gotCity.South.Name,
				"city %q south direction mismatch", name,
			)
		} else {
			require.Nil(t, gotCity.South, "city %q south direction not nil", name)
		}
	}
}
