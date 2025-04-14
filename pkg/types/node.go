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

// NOTE: must improve
var RetType = map[string]ModoType{
	OPERATOR_ADD:   TY_INT32,
	OPERATOR_EQ:    TY_BOOL,
	OPERATOR_GT:    TY_BOOL,
	OPERATOR_LT:    TY_BOOL,
	OPERATOR_AND:   TY_BOOL,
	OPERATOR_OR:    TY_BOOL,
	LIB_CORE_CONJ:  TY_VECTOR,
	LIB_CORE_ASSOC: TY_VECTOR,
	LIB_CORE_POP:   TY_VECTOR,
	LIB_CORE_GET:   TY_INT32, // HACK: is not TRUE
}

type Program struct {
	Declares    *Node
	BuiltinLibs *BuiltinLibProp
	GlobalStr   []*ir.InstLoad
}

type BuiltinLibProp struct {
	GlobalVar *BuiltinGlobalVarsProp
	Printf    *BuiltinProp
	Malloc    *BuiltinProp
	Memcpy    *BuiltinProp
}

// HACK: should remove llir/ir reference from this namespace
type BuiltinProp struct {
	FuncPtr *ir.Func
}

// HACK: should remove llir/ir reference from this namespace
type BuiltinGlobalVarsProp struct {
	FormatDigit        *ir.Global
	FormatStr          *ir.Global
	FormatSpace        *ir.Global
	FormatCR           *ir.Global
	TrueValue          *ir.Global
	FalseValue         *ir.Global
	NilValue           *ir.Global
	FormatBracketOpen  *ir.Global
	FormatBracketClose *ir.Global
	FormatComma        *ir.Global
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
	VarPtr   value.Value // binded local variable
	FuncPtr  *ir.Func    // declared function, library function
	IRValue  value.Value //
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

func (node *Node) GetBlockName(s string) string {
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
func GetLLVMTypeForVector(node *Node) (types.Type, bool) {

	elemType, ok := GetLLVMType(node.ElemType)
	if !ok {
		return nil, false
	}
	var length uint64 = 0
	for n := node.Child; n != nil; n = n.Next {
		length++
	}

	var typeMap = map[ModoType]types.Type{
		TY_VECTOR: &types.PointerType{
			ElemType: &types.ArrayType{
				ElemType: elemType,
				Len:      length,
			},
		},
	}

	if t, ok := typeMap[node.Type]; ok {
		return t, true
	}
	return nil, true
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
	log.Debug(
		log.BLUE(
			fmt.Sprintf(
				"%s %p %#+v %#+v %#+v %d %#+v",
				strings.Repeat("  ", depth),
				node,
				node.Kind,
				node.Type,
				node.Val,
				node.Len,
				node.ElemType,
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
	if prog.BuiltinLibs != nil {
		log.DebugMessage("[BuiltinLib]")
	}
}
