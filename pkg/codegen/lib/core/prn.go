package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PrnSchalar(
	formatStr *ir.Global,
	libs *mTypes.BuiltinLibProp,
	block *ir.Block,
	n *mTypes.Node,
) {
	value := n.IRValue
	ty := n.IRValue.Type()
	if ty.Equal(types.I32) {
		formatStr = libs.GlobalVar.FormatDigit

	} else if ty.Equal(types.I1) {
		formatStr = libs.GlobalVar.FormatStr
		value = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)

	} else if ty.Equal(types.I8Ptr) {
		formatStr = libs.GlobalVar.FormatStr

	} else if ty.Equal(types.Void) {
		formatStr = libs.GlobalVar.FormatStr
		value = libs.GlobalVar.NilValue
	}
	block.NewCall(libs.Printf.FuncPtr, formatStr, value)
}

func PrnVector(libs *mTypes.BuiltinLibProp, block *ir.Block, n *mTypes.Node, t types.Type) {
	value, ok := n.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", value)
	}

	block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatBracketOpen)
	for i := uint64(0); i < uint64(n.Len); i++ {
		elemPtr := block.NewGetElementPtr(
			value.ElemType,
			value,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		elem := block.NewLoad(types.I32, elemPtr)
		// NOTE: format may be changable
		block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatDigit, elem)

		if i < uint64(n.Len-1) {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatComma)
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}
	}
	block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatBracketClose)
}

func InvokePrn(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var formatStr *ir.Global

	for n := node; n != nil; n = n.Next {
		ty := n.IRValue.Type()

		if ty.Equal(types.I32) ||
			ty.Equal(types.I1) ||
			ty.Equal(types.I8Ptr) ||
			ty.Equal(types.Void) {
			PrnSchalar(formatStr, libs, block, n)

		} else if t, ok := ty.(*types.PointerType); ok {
			PrnVector(libs, block, n, t)

		} else {
			log.Panic("unresolved type: have %+v", n)
		}

		if n.Next == nil {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatCR)
		} else {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}

	}
	return nil
}
