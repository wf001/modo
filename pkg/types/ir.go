package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
)

func GetPrintFormat(ty types.Type, libs *BuiltinLibProp) (*ir.Global, bool) {

	formatMap := map[types.Type]*ir.Global{
		types.I1:    libs.GlobalVar.FormatStr,
		types.I8Ptr: libs.GlobalVar.FormatStr,
		types.I32:   libs.GlobalVar.FormatDigit,
		types.Void:  libs.GlobalVar.FormatStr,
	}
	if f, ok := formatMap[ty]; ok {
		return f, true
	}

	log.Debug("unresolved type: have %+v", ty)
	return nil, false
}

func IsScalar(v value.Value) bool {
	return v.Type().Equal(types.I1) ||
		v.Type().Equal((types.I1Ptr)) ||
		v.Type().Equal(types.I8Ptr) ||
		v.Type().Equal(types.I32) ||
		v.Type().Equal((types.I32Ptr)) ||
		v.Type().Equal(types.Void)
}

func LoadArrElem(block *ir.Block, src value.Value, ty *types.ArrayType, i uint64) *ir.InstLoad {
	elemPtr := block.NewGetElementPtr(
		ty,
		src,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(i)),
	)
	return block.NewLoad(ty.ElemType, elemPtr)

}

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
