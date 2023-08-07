package mapfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Structure struct {
	Cities Cities
}

// ReadFile reads map from file name
func ReadFile(fileName string) (*Structure, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return Read(f)
}

// Read reads map from a readable stream
func Read(r io.Reader) (*Structure, error) {
	cities := make(map[string]*Node)

	lineNo := 1
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if str == "" {
			continue
		}

		err := parseLine(str, cities)
		if err != nil {
			return nil, fmt.Errorf("%w (line: %d)", err, lineNo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Structure{
		Cities: cities,
	}, nil
}

// parseLine parses single map line and appends/updates city in passed map
func parseLine(l string, dst map[string]*Node) error {
	chunks := strings.Split(l, " ")
	cityName := chunks[0]

	cityNode, ok := dst[cityName]
	if !ok {
		cityNode = &Node{Name: cityName}
		dst[cityName] = cityNode
	}

	hasDirections := false
	for _, opts := range chunks[1:] {
		opts = strings.TrimSpace(opts)
		if opts == "" {
			// Skip double whitespace
			continue
		}

		hasDirections = true
		fields := strings.Split(opts, "=")
		if len(fields) != 2 || fields[1] == "" {
			return fmt.Errorf("invalid direction attribute: %q", opts)
		}

		key, value := fields[0], fields[1]
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
			return fmt.Errorf("invalid city direction: %q", opts)
		}
	}

	if !hasDirections {
		return fmt.Errorf("missing directions for city %q", cityName)
	}

	return nil
}
