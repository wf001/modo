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
	ND_PROGRAM_ROOT = NodeKind("nd.program.root")

	// reserved symbol
	ND_VAR_DECLARE   = NodeKind("nd.var.declare")   // variables and functions
	ND_VAR_REFERENCE = NodeKind("nd.var.reference") // variables and functions
	ND_TYPE_DECLARE  = NodeKind("nd.type.declare")
	ND_DECLARE       = NodeKind("nd.declare")  // def
	ND_LAMBDA        = NodeKind("nd.lambda")   // fn
	ND_BIND          = NodeKind("nd.bind")     // let
	ND_EXPR          = NodeKind("nd.expr")     // set of functions
	ND_IF            = NodeKind("nd.if")       // if
	ND_FUNCCALL      = NodeKind("nd.funccall") // calling declared function
	ND_LIBCALL       = NodeKind("nd.libcall")  // calling builtin library function

	// type
	ND_SCALAR     = NodeKind("nd.scalar")     // int32, string, bool
	ND_COLLECTION = NodeKind("nd.collection") // vector, map, list
)

type ModoType string

// Note: what for?
const (
	TY_INT32    = ModoType("ty.int32")
	TY_STR      = ModoType("ty.str")
	TY_NIL      = ModoType("ty.nil")
	TY_BOOL     = ModoType("ty.bool")
	TY_VECTOR   = ModoType("ty.vector")
	TY_STRUCT   = ModoType("ty.struct")
	TY_EXTENDED = ModoType("ty.extended")
)

type NodeType struct {
	Value      ModoType
	ExtendName string
	Child      *NodeType
}

func (n *NodeType) String() string {

	if n.Child != nil {
		return fmt.Sprintf(
			"{Value: \"%s\", ExtendName: \"%s\", Child: \"%s\" }",
			n.Value,
			n.ExtendName,
			n.Child.String(),
		)
	}

	return fmt.Sprintf(
		"{Value: \"%s\", ExtendName: \"%s\" }",
		n.Value,
		n.ExtendName,
	)
}

type Node struct {
	Kind     NodeKind
	Next     *Node
	Type     *NodeType
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
		"{Kind:%#+v, Val:%#+v, Type:%s,  Len:%d}",
		n.Kind,
		n.Val,
		n.Type,
		n.Len,
	)
}

// ==============
// predication
// ==============

func (node *Node) IsKind(kind NodeKind) bool {
	return node.Kind == kind
}

func (node *Node) IsType(ty ModoType) bool {
	return node.Type.Value == ty
}

func (node *Node) IsScalar() bool {
	if node.Type == nil {
		return false
	}
	return node.IsType(TY_INT32) ||
		node.IsType(TY_BOOL) ||
		node.IsType(TY_STR) ||
		node.IsType(TY_NIL)
}

// ==============
// naming
// ==============

func GetGlobalVarName(s string, m *ir.Module, ndtype *NodeType) string {
	return fmt.Sprintf(".%s.%p", s, ndtype)
}

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

func isStructTypeDefined(mod *ir.Module, fields []types.Type) (types.Type, bool) {
	for _, typ := range mod.TypeDefs {
		structType, ok := typ.(*types.StructType)
		if !ok {
			continue
		}
		// 明示的に同じかチェック（要素数・型一致）
		if len(structType.Fields) != len(fields) {
			continue
		}
		match := true
		for i := range fields {
			if !types.Equal(structType.Fields[i], fields[i]) {
				match = false
				break
			}
		}
		if match {
			return structType, true
		}
	}
	return nil, false
}

func DeclareVectorType(ir *ir.Module, ty types.Type, ndtype *NodeType) types.Type {
	// i32 type vector
	// 型: struct { i32* %arrElm, i64 %len}
	vectorIntType := types.NewStruct(types.NewPointer(ty), types.I64)

	// not declare same type twice
	if ty, defined := isStructTypeDefined(ir, vectorIntType.Fields); defined {
		return ty
	}
	typeName := GetGlobalVarName("vec", ir, ndtype)
	vectorIntType.SetName(typeName)
	ir.NewTypeDef(typeName, vectorIntType)
	return vectorIntType

}

// Get LLVM type from corresponding Node Type
func GetLLVMTypeRec(
	m *ir.Module,
	ndtype *NodeType,
	prelude *PreludeProps,
) (types.Type, types.Type, bool) {
	var scalarTy types.Type

	var scalarTypeMap = map[ModoType]types.Type{
		TY_INT32: types.I32,
		TY_BOOL:  types.I1,
		TY_STR:   types.I8Ptr,
		TY_NIL:   types.Void,
	}

	scalarTy, isRootScalar := scalarTypeMap[ndtype.Value]
	if isRootScalar {
		return scalarTy, nil, true
	} else if ndtype.Child == nil {
		return nil, nil, false
	}

	if ndtype.Value == TY_VECTOR {
		chidType, _, _ := GetLLVMTypeRec(m, ndtype.Child, prelude)
		ret := DeclareVectorType(m, chidType, ndtype)
		return ret, chidType, true
	}

	return nil, nil, false
}
func GetExtendedType(declare DeclareProps, node *Node) *StructType {
	for k, v := range declare.Type.Struct {
		if k == node.Type.ExtendName {
			return v
		}
	}
	return nil
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
	if prog.Declare.Func != nil {
		log.DebugMessage("[Declares]")
		prog.Declare.Func.debugRecursive(0)
	}

}
