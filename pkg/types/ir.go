package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func GetPrintFormat(ty types.Type, libs *BuiltinLibProp) (*ir.Global, bool) {

	formatMap := map[types.Type]*ir.Global{
		types.I32:   libs.GlobalVar.FormatDigit,
		types.I1:    libs.GlobalVar.FormatStr,
		types.I8Ptr: libs.GlobalVar.FormatStr,
		types.Void:  libs.GlobalVar.FormatStr,
	}
	if f, ok := formatMap[ty]; ok {
		return f, true
	}

	return nil, false
}

func IsScalar(v value.Value) bool {
	return v.Type().Equal(types.I32) ||
		v.Type().Equal(types.I1) ||
		v.Type().Equal(types.I8Ptr) ||
		v.Type().Equal(types.Void)

}
