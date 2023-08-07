package mapfile

const DirectionsCount = 4

// Cities is key-value pair of city name and city node.
type Cities map[string]*Node

// Names returns a slice of city names.
func (cities Cities) Names() []string {
	names := make([]string, 0, len(cities))
	for name := range cities {
		names = append(names, name)
	}

	return names
}

// AsSlice returns cities as slice.
func (cities Cities) AsSlice() []*Node {
	slice := make([]*Node, 0, len(cities))
	for _, elem := range cities {
		slice = append(slice, elem)
	}

	return slice
}

// Node represents a city node in graph and its siblings.
type Node struct {
	Name  string
	South *Node
	North *Node
	East  *Node
	West  *Node
}

// Directions returns a slice of available valid directions from city.
func (n *Node) Directions() []*Node {
	directions := make([]*Node, 0, DirectionsCount)
	if n.East != nil {
		directions = append(directions, n.East)
	}
	if n.West != nil {
		directions = append(directions, n.West)
	}
	if n.South != nil {
		directions = append(directions, n.South)
	}
	if n.North != nil {
		directions = append(directions, n.North)
	}

	return directions
}

// String implements fmt.Stringer.
func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.Name
}
