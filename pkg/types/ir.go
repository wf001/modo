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
	Struct map[string]*PreludeStruct
	LLVM   map[string]*types.Type
}

type PreludeProps struct {
	Types PreludeTypeProps
}

type PreludeTypeProps struct {
	VectorInt    *types.StructType
	VectorString *types.StructType
	VectorBool   *types.StructType
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

type PreludeStruct struct {
	Name  string
	Field map[string]PreludeStructFields
	Types *types.StructType
}

type PreludeStructFields struct {
	Pos  uint64
	Type types.Type
}

// ==============
// predication
// ==============

func IsScalar(v value.Value) bool {
	return v.Type().Equal(types.I1) ||
		v.Type().Equal(types.I8Ptr) ||
		v.Type().Equal(types.I32) ||
		v.Type().Equal(types.Void)
}

// ==============
// conversion llir/llvm properties to other properties
// ==============

func GetPrintFormat(ty types.Type, internal *Internal) (*ir.Global, bool) {

	formatMap := map[types.Type]*ir.Global{
		types.I1:    internal.GlobalConst.FormatStr,
		types.I8Ptr: internal.GlobalConst.FormatStr,
		types.I32:   internal.GlobalConst.FormatDigit,
		types.Void:  internal.GlobalConst.FormatStr,
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
	case *types.PointerType:
		// 通常は64bitだが、プラットフォームによって異なる
		return 64
	default:
		return 0 // 未対応型など
	}
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func LoadArrElem(block *ir.Block, src value.Value, ty *types.ArrayType, i uint64) *ir.InstLoad {
	elemPtr := block.NewGetElementPtr(
		ty,
		src,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(i)),
	)
	return block.NewLoad(ty.ElemType, elemPtr)

}

// Note: remain here until it will be defined the strategy of memory lifecycle
func AssertArrType(v value.Value) (*types.ArrayType, bool) {
	ptrElemType, ok := v.Type().(*types.PointerType)

	if !ok {
		log.Warn("Not pointer type, got: %+v", v.Type())
		return nil, false

	}
	arrType, ok := ptrElemType.ElemType.(*types.ArrayType)

	if !ok {
		log.Warn("Not pointer type to array, got: %+v", ptrElemType.ElemType)
		return nil, false
	}
	return arrType, true
}
