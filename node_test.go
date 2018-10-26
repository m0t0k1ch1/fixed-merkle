package merkle

import (
	"bytes"
	"testing"
)

func testNodesEquality(t *testing.T, expected, actual *Node) {
	if !bytes.Equal(actual.Bytes(), expected.Bytes()) {
		t.Errorf("expected: %x, actual: %x", expected.Bytes(), actual.Bytes())
	}
	if actual.String() != expected.String() {
		t.Errorf("expected: %s, actual: %s", expected.String(), actual.String())
	}
}
