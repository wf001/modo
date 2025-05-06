package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	node := &Node{Val: "testFunc"}
	want := "fn.testFunc"
	assert.Equal(t, want, node.GetFuncName())
}

func TestGetFuncNameUnnamed(t *testing.T) {
	node := &Node{}
	want := "fn.unnamed." + fmt.Sprintf("%p", node)
	assert.Equal(t, want, node.GetUnnamedFuncName())
}

func add(kind NodeKind, val string) *Node {
	return &Node{
		Kind: kind,
		Val:  val,
	}
}

func buildNode(tokens []*Node) *Node {
	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Next = tokens[i+1]
	}
	return tokens[0]
}

func TestGetLastNode(t *testing.T) {
	node := []*Node{
		add(ND_PROGRAM_ROOT, ""),
		add(ND_LIBCALL, ""),
		add(ND_SCALAR, ""),
	}
	want := add(ND_SCALAR, "")
	assert.Equal(t, want, buildNode(node).GetLastNode())
}

func TestGetNodeSize(t *testing.T) {
	childNode := add(ND_SCALAR, "")
	childNode.Next = add(ND_SCALAR, "")
	childNode.Next.Next = add(ND_SCALAR, "")

	want := uint64(3)
	assert.Equal(t, want, childNode.GetNodeSize())
}
