package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokePrn(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var formatStr *ir.Global

	for n := node; n != nil; n = n.Next {
		value := n.IRValue
		ty := n.IRValue.Type()

		if ty.Equal(types.I32) {
			formatStr = libs.GlobalVar.FormatDigit

		} else if ty.Equal(types.I1) {
			formatStr = libs.GlobalVar.FormatStr
			value = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)

		} else if ty.Equal(types.I8Ptr) {
			formatStr = libs.GlobalVar.FormatStr

		} else if _, ok := ty.(*types.VectorType); ok {
			formatStr = libs.GlobalVar.FormatDigit
			elem := block.NewExtractElement(value, constant.NewInt(types.I32, int64(1)))
			value = block.NewAdd(constant.NewInt(types.I32, 0), elem)

		} else if ty.Equal(types.Void) {
			formatStr = libs.GlobalVar.FormatStr
			value = libs.GlobalVar.NilValue

		} else {
			log.Panic("unresolved type: have %+v", n)
		}
		block.NewCall(libs.Printf.FuncPtr, formatStr, value)

		if n.Next == nil {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatCR)
		} else {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}

	}
	return nil
}
