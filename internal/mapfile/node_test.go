package mapfile

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_Directions(t *testing.T) {
	cases := map[string]struct {
		node *Node
		want []*Node
	}{
		"node with all directions": {
			node: &Node{
				East:  &Node{Name: "East"},
				West:  &Node{Name: "West"},
				South: &Node{Name: "South"},
				North: &Node{Name: "North"},
			},
			want: []*Node{
				{Name: "East"},
				{Name: "West"},
				{Name: "South"},
				{Name: "North"},
			},
		},
		"node with partial directions": {
			node: &Node{
				East: &Node{Name: "East"},
			},
			want: []*Node{
				{Name: "East"},
			},
		},
		"node without directions": {
			node: &Node{},
			want: []*Node{},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			got := c.node.Directions()
			require.Equal(t, c.want, got)
		})
	}
}

func TestNode_String(t *testing.T) {
	cases := map[string]struct {
		node *Node
		want string
	}{
		"nil": {
			node: nil,
			want: "<nil>",
		},
		"not-nil": {
			node: &Node{Name: "foo"},
			want: "foo",
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			got := c.node.String()
			require.Equal(t, c.want, got)
		})
	}
}

func TestCities_Names(t *testing.T) {
	cities := Nodes{
		"foo": {Name: "foo"},
		"bar": {Name: "Bar"},
	}

	want := []string{"foo", "bar"}
	got := cities.Names()

	sort.Strings(want)
	sort.Strings(got)
	require.Equal(t, want, got)
}

func TestCities_AsSlice(t *testing.T) {
	cities := Nodes{
		"foo": {Name: "foo"},
	}

	expect := []*Node{
		{Name: "foo"},
	}

	got := cities.AsSlice()
	require.Equal(t, expect, got)
}
