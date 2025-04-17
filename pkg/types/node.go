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

type Program struct {
	Declares          *Node
	Prelude           *PreludeProps
	DeclaredGlobalVar []*ir.InstLoad
	Internal          *Internal
}

type PreludeProps struct {
	Types PreludeTypeProps
}

type PreludeTypeProps struct {
	VectorInt    *types.StructType
	VectorString *types.StructType
}

type Internal struct {
	Cstd        *Cstd
	GlobalConst *GlobalConst
}

type Cstd struct {
	Printf *ir.Func
	Malloc *ir.Func
	Memcpy *ir.Func
	Strcmp *ir.Func
}

// HACK: should remove llir/ir reference from this namespace
type GlobalConst struct {
	FormatDigit        *ir.Global
	FormatStr          *ir.Global
	StringSpace        *ir.Global
	StringCR           *ir.Global
	StringTrue         *ir.Global
	StringFalse        *ir.Global
	StringNil          *ir.Global
	StringBracketOpen  *ir.Global
	StringBracketClose *ir.Global
	StringComma        *ir.Global
}

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

// pred kind
func (node *Node) IsKind(kind NodeKind) bool {
	return node.Kind == kind
}

// pred type
func (node *Node) IsType(ty ModoType) bool {
	return node.Type == ty
}

func (node *Node) IsScalar() bool {
	return node.IsType(TY_INT32) ||
		node.IsType(TY_BOOL) ||
		node.IsType(TY_STR) ||
		node.IsType(TY_NIL)
}

// naming
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

// Returns the last node of the linked list.
func (node *Node) GetLastNode() *Node {
	lastNode := node
	for lastNode.Next != nil {
		lastNode = lastNode.Next
	}
	return lastNode
}

// Get LLVM type from corresponding custom type
func GetLLVMType(ty ModoType) (types.Type, bool) {

	var typeMap = map[ModoType]types.Type{
		TY_INT32: types.I32,
		TY_BOOL:  types.I1,
		TY_STR:   types.I8Ptr,
		TY_NIL:   types.Void,
	}

	if t, ok := typeMap[ty]; ok {
		return t, true
	}
	return nil, false
}
func GetNodeSize(node *Node) uint64 {
	var length uint64 = 0
	for n := node.Child; n != nil; n = n.Next {
		length++
	}
	return length

}
func GetLLVMTypeForVector(node *Node, prelude *PreludeProps) (types.Type, bool) {

	var typeMap = map[ModoType]types.Type{
		TY_INT32: prelude.Types.VectorInt,
		TY_STR:   prelude.Types.VectorString,
	}
	if node.Type != TY_VECTOR {
		return nil, true
	}

	if t, ok := typeMap[node.ElemType]; ok {
		return t, true
	}
	return nil, true
}

func GetBitWidth(t types.Type) uint64 {
	switch typ := t.(type) {
	case *types.IntType:
		return typ.BitSize
	case *types.FloatType:
		switch typ.Kind {
		case types.FloatKindHalf:
			return 16
		case types.FloatKindFloat:
			return 32
		case types.FloatKindDouble:
			return 64
		default:
			return 0
		}
	case *types.PointerType:
		// 通常は64bitだが、プラットフォームによって異なる
		return 64
	default:
		return 0 // 未対応型など
	}
}

// debug
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

func (node *Node) Debug(depth int) {
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
		node.Bind.Debug(depth + 1)

		indicate(".Child", depth+1)
		node.Child.Debug(depth + 1)
	case ND_LAMBDA:
		indicate(".Args", depth+1)
		node.Args.Debug(depth + 1)

		indicate(".Child", depth+1)
		node.Child.Debug(depth + 1)
	case ND_IF:
		indicate(".Cond", depth+1)
		node.Cond.Debug(depth + 1)

		indicate(".Then", depth+1)
		node.Then.Debug(depth + 1)

		indicate(".Else", depth+1)
		node.Else.Debug(depth + 1)
	default:
		node.Child.Debug(depth + 1)
	}
	if node.Next != nil {
		node.Next.Debug(depth)
	}
}
func (prog *Program) Debug(depth int) {
	if prog.Declares != nil {
		log.DebugMessage("[Declares]")
		prog.Declares.Debug(0)
	}

	// TODO: implement
	// NOTE: is type of BuiltinLibs Node?
}
