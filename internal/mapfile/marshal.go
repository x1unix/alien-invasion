package mapfile

import "strings"

// MarshalNodes encodes nodes as map file string.
func MarshalNodes(nodes Nodes) string {
	if len(nodes) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for name, node := range nodes {
		sb.WriteString(name)
		if node.East != nil {
			sb.WriteString(" east=")
			sb.WriteString(node.East.Name)
		}
		if node.West != nil {
			sb.WriteString(" west=")
			sb.WriteString(node.West.Name)
		}
		if node.South != nil {
			sb.WriteString(" south=")
			sb.WriteString(node.South.Name)
		}
		if node.North != nil {
			sb.WriteString(" north=")
			sb.WriteString(node.North.Name)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
