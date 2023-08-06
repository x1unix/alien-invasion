package mapfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const DirectionsCount = 4

type Cities map[string]*Node

func (cities Cities) AsSlice() []*Node {
	slice := make([]*Node, 0, len(cities))
	for _, elem := range cities {
		slice = append(slice, elem)
	}

	return slice
}

type Structure struct {
	Cities   Cities
	RootCity *Node
}

type Node struct {
	Name                     string
	South, North, East, West *Node
}

func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.Name
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

func ReadFile(fileName string) (*Structure, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return Read(f)
}

func Read(r io.Reader) (*Structure, error) {
	var rootNode *Node
	cities := make(map[string]*Node)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if str == "" {
			continue
		}

		cityNode, err := parseLine(str, cities)
		if err != nil {
			return nil, err
		}

		if rootNode == nil {
			rootNode = cityNode
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Structure{
		RootCity: rootNode,
		Cities:   cities,
	}, nil
}

func parseLine(l string, dst map[string]*Node) (*Node, error) {
	chunks := strings.Split(l, " ")
	cityName := chunks[0]

	cityNode, ok := dst[cityName]
	if !ok {
		cityNode = &Node{Name: cityName}
		dst[cityName] = cityNode
	}

	for _, opts := range chunks[1:] {
		opts = strings.TrimSpace(opts)
		if len(opts) == 0 {
			continue
		}

		fields := strings.Split(opts, "=")
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid line: %q", opts)
		}

		key := fields[0]
		value := fields[1]

		sibling, ok := dst[value]
		if !ok {
			sibling = &Node{Name: value}
			dst[value] = sibling
		}

		switch key {
		case "west":
			cityNode.West = sibling
		case "east":
			cityNode.East = sibling
		case "north":
			cityNode.North = sibling
		case "south":
			cityNode.South = sibling
		default:
			return nil, fmt.Errorf("invalid direction in line: %q", opts)
		}
	}

	return cityNode, nil
}
