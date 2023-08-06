package mapfile

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadFile(t *testing.T) {
	node, err := ReadFile("testdata/file.map")
	require.NoError(t, err)

	t.Log(node)
}
