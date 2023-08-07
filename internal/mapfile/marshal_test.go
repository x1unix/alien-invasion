package mapfile

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalNodes(t *testing.T) {
	cases := map[string]struct {
		want  []string
		input Nodes
	}{
		"should skip if when no nodes": {
			want:  []string{""},
			input: nil,
		},
		"should print node with siblings": {
			want: []string{
				"Whiterun east=Windhelm west=Markarth south=Falkreath north=Dawnstar",
				"Riften west=Falkreath north=Windhelm",
				"Windhelm west=Whiterun south=Riften",
				"Falkreath east=Riften north=Whiterun",
				"Markarth east=Whiterun",
				"Dawnstar south=Whiterun",
			},
			input: Nodes{
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
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			output := MarshalNodes(c.input)
			got := strings.Split(strings.TrimSpace(output), "\n")

			sort.Strings(got)
			sort.Strings(c.want)
			require.Equal(t, c.want, got)
		})
	}
}
