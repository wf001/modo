package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func GetPrintFormat(ty types.Type, libs *BuiltinLibProp) *ir.Global {

	if ty.Equal(types.I32) {
		return libs.GlobalVar.FormatDigit

	} else if ty.Equal(types.I1) {
		return libs.GlobalVar.FormatStr

	} else if ty.Equal(types.I8Ptr) {
		return libs.GlobalVar.FormatStr

	} else if ty.Equal(types.Void) {
		return libs.GlobalVar.FormatStr
	}
	return nil
}

// Note: remove either this or codegen.isConstant
func IsScalar(v value.Value) bool {
	return v.Type().Equal(types.I32) ||
		v.Type().Equal(types.I1) ||
		v.Type().Equal(types.I8Ptr) ||
		v.Type().Equal(types.Void)

}
