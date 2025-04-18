package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
)

type NodeKind string

const (
	ND_PROGRAM_ROOT = NodeKind("ND_PROGRAM_ROOT")

	// reserved symbol
	ND_VAR_DECLARE   = NodeKind("ND_VAR_DECLARE")   // variables and functions
	ND_VAR_REFERENCE = NodeKind("ND_VAR_REFERENCE") // variables and functions
	ND_DECLARE       = NodeKind("ND_DECLARE")       // def
	ND_LAMBDA        = NodeKind("ND_LAMBDA")        // fn
	ND_BIND          = NodeKind("ND_BIND")          // let
	ND_EXPR          = NodeKind("ND_EXPR")          // set of functions
	ND_IF            = NodeKind("ND_IF")            // if
	ND_FUNCCALL      = NodeKind("ND_FUNCCALL")
	ND_LIBCALL       = NodeKind("ND_LIBCALL")

	// type
	ND_SCALAR     = NodeKind("ND_SCALAR")     // int32, string, bool
	ND_COLLECTION = NodeKind("ND_COLLECTION") // vector, map, list
)

type ModoType string

// Note: what for?
const (
	TY_INT32  = ModoType("TY_INT32")
	TY_STR    = ModoType("TY_STR")
	TY_NIL    = ModoType("TY_NIL")
	TY_BOOL   = ModoType("TY_BOOL")
	TY_VECTOR = ModoType("TY_VECTOR")
)

type Node struct {
	Kind     NodeKind
	Next     *Node
	Type     ModoType
	ElemType ModoType
	Child    *Node
	Cond     *Node
	CondRet  *ir.InstAlloca
	Then     *Node
	Else     *Node
	Val      string
	Len      uint64 // the number of bytes, used with string type
	Bind     *Node
	Args     *Node
	IsGlobal bool
	FuncPtr  *ir.Func    // declared function, library function
	IRValue  value.Value //
}

func (n *Node) String() string {

	return fmt.Sprintf(
		"{Kind:%#+v, Val:%#+v, Type:%#+v, Len:%d, ElemType:%#+v}",
		n.Kind,
		n.Val,
		n.Type,
		n.Len,
		n.ElemType,
	)
}

// ==============
// predication
// ==============

func (node *Node) IsKind(kind NodeKind) bool {
	return node.Kind == kind
}

func (node *Node) IsType(ty ModoType) bool {
	return node.Type == ty
}

func (node *Node) IsScalar() bool {
	return node.IsType(TY_INT32) ||
		node.IsType(TY_BOOL) ||
		node.IsType(TY_STR) ||
		node.IsType(TY_NIL)
}

// ==============
// naming
// ==============

func (node *Node) GetUnnamedFuncName() string {
	return fmt.Sprintf("fn.%s.%p", "unnamed", node)
}

func (node *Node) GetFuncName() string {
	return fmt.Sprintf("fn.%s", node.Val)
}

func (node *Node) GetBlockName(s string, blks []*ir.Block) string {
	return fmt.Sprintf("%s.%p.%d", s, node, len(blks))
}

func (node *Node) GetVarName(s string, insts []ir.Instruction) string {
	return fmt.Sprintf("%s.%p", s, node)
}

// ==============
// knowing the properties of Node
// ==============

// Returns the last node of the linked list.
func (node *Node) GetLastNode() *Node {
	lastNode := node
	for lastNode.Next != nil {
		lastNode = lastNode.Next
	}
	return lastNode
}

func (node *Node) GetNodeSize() uint64 {
	var size uint64 = 0
	for n := node.Child; n != nil; n = n.Next {
		size++
	}
	return size

}

// ==============
// conversion Node properties to other properties
// ==============

// Get LLVM type from corresponding Node Type
func GetLLVMType(node *Node, prelude *PreludeProps) (types.Type, types.Type, bool) {
	var scalarTy types.Type

	var scalarTypeMap = map[ModoType]types.Type{
		TY_INT32: types.I32,
		TY_BOOL:  types.I1,
		TY_STR:   types.I8Ptr,
		TY_NIL:   types.Void,
	}

	scalarTy, isRootScalar := scalarTypeMap[node.Type]
	if isRootScalar {
		return nil, scalarTy, true
	}

	var collTy types.Type
	var collTypeMap = map[ModoType]types.Type{
		TY_INT32: prelude.Types.VectorInt,
		TY_STR:   prelude.Types.VectorString,
	}

	collTy, isRootColl := collTypeMap[node.ElemType]
	scalarTy, isElemScalar := scalarTypeMap[node.ElemType]

	if node.Type == TY_VECTOR {
		if isRootColl && isElemScalar {
			return collTy, scalarTy, true
		}
	}

	return nil, nil, false
}

// ==============
// The following is for developement purposes
// ==============

func indicate(s string, depth int) {
	log.Debug(
		log.YELLOW(
			fmt.Sprintf(
				"%s [%s]",
				strings.Repeat("  ", depth),
				s,
			),
		))

}

func (node *Node) debugRecursive(depth int) {
	if node == nil {
		return
	}
	log.DebugNoLine(
		log.BLUE(
			fmt.Sprintf(
				"%s %s",
				strings.Repeat("  ", depth),
				node,
			),
		),
	)

	switch node.Kind {
	case ND_BIND:
		indicate(".Bind", depth+1)
		node.Bind.debugRecursive(depth + 1)

		indicate(".Child", depth+1)
		node.Child.debugRecursive(depth + 1)
	case ND_LAMBDA:
		indicate(".Args", depth+1)
		node.Args.debugRecursive(depth + 1)

		indicate(".Child", depth+1)
		node.Child.debugRecursive(depth + 1)
	case ND_IF:
		indicate(".Cond", depth+1)
		node.Cond.debugRecursive(depth + 1)

		indicate(".Then", depth+1)
		node.Then.debugRecursive(depth + 1)

		indicate(".Else", depth+1)
		node.Else.debugRecursive(depth + 1)
	default:
		node.Child.debugRecursive(depth + 1)
	}
	if node.Next != nil {
		node.Next.debugRecursive(depth)
	}
}

func (prog *Program) Debug(depth int) {
	if prog.Declares != nil {
		log.DebugMessage("[Declares]")
		prog.Declares.debugRecursive(0)
	}

}
