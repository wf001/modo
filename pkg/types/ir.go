package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
)

type Program struct {
	Declare DeclareProps
	Prelude *PreludeProps
}

type DeclareProps struct {
	Func      *Node
	GlobalVar []*ir.InstLoad
	Type      *ExtendedTypes
}

type ExtendedTypes struct {
	Struct map[string]*StructType
	LLVM   map[string]*types.Type
}

type PreludeProps struct {
	Types PreludeTypeProps
}

type PreludeTypeProps struct {
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

type GlobalConst struct {
	FormatDigit        *ir.Global
	FormatDouble       *ir.Global
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

type StructType struct {
	Name  string
	Field map[string]StructTypeField
	Types *types.StructType
}

type StructTypeField struct {
	Pos  uint64
	Type types.Type
}

// ==============
// constant value
// ==============
var I1zero = constant.NewInt(types.I1, 0)
var I1one = constant.NewInt(types.I1, 1)

var I32zero = constant.NewInt(types.I32, 0)
var I32one = constant.NewInt(types.I32, 1)

var I64zero = constant.NewInt(types.I64, 0)
var I64one = constant.NewInt(types.I64, 1)

// ==============
// predication
// ==============

func IsScalar(v value.Value) bool {
	return v.Type().Equal(types.I1) ||
		v.Type().Equal(types.I8Ptr) ||
		v.Type().Equal(types.I32) ||
		v.Type().Equal(types.Double) ||
		v.Type().Equal(types.Void)
}

func TypeExists(declare DeclareProps, typeName string) bool {

	for k := range declare.Type.Struct {
		if k == typeName {
			return true
		}
	}
	for k := range declare.Type.LLVM {
		if k == typeName {
			return true
		}
	}
	return false
}

// ==============
// conversion llir/llvm properties to other properties
// ==============

func GetPrintFormat(ty types.Type, internal *Internal) (*ir.Global, bool) {

	formatMap := map[types.Type]*ir.Global{
		types.I1:     internal.GlobalConst.FormatStr,
		types.I8Ptr:  internal.GlobalConst.FormatStr,
		types.I32:    internal.GlobalConst.FormatDigit,
		types.Double: internal.GlobalConst.FormatDouble,
		types.Void:   internal.GlobalConst.FormatStr,
	}
	if f, ok := formatMap[ty]; ok {
		return f, true
	}

	log.Debug("unresolved type: have %+v", ty)
	return nil, false
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
	case *types.PointerType, *types.StructType:
		// 通常は64bitだが、プラットフォームによって異なる
		return 64
	default:
		return 0 // 未対応型など
	}
}

func GetVectorTypeFromPtr(v value.Value) (*types.StructType, types.Type) {
	structPtr := v.Type().(*types.PointerType)
	structedVecType := structPtr.ElemType.(*types.StructType)
	elemType := structedVecType.Fields[0].(*types.PointerType).ElemType

	return structedVecType, elemType
}
