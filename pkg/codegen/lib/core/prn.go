package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func prnScalar(
	libs *mTypes.BuiltinLibProp,
	block *ir.Block,
	n *mTypes.Node,
) {
	value := n.IRValue
	ty := n.IRValue.Type()
	formatStr := mTypes.GetPrintFormat(ty, libs)

	if ty.Equal(types.I1) {
		value = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)
	} else if ty.Equal(types.Void) {
		value = libs.GlobalVar.NilValue
	}

	block.NewCall(libs.Printf.FuncPtr, formatStr, value)
}

func prnVector(libs *mTypes.BuiltinLibProp, block *ir.Block, n *mTypes.Node) {
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
		formatStr := mTypes.GetPrintFormat(elem.ElemType, libs)
		// NOTE: format may be changable
		block.NewCall(libs.Printf.FuncPtr, formatStr, elem)

		if i < uint64(n.Len-1) {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatComma)
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}
	}

	block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatBracketClose)
}

func InvokePrn(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {

	for n := node; n != nil; n = n.Next {
		ty := n.IRValue.Type()

		if mTypes.IsScalar(n.IRValue) {
			prnScalar(libs, block, n)

		} else if _, ok := ty.(*types.PointerType); ok {
			prnVector(libs, block, n)

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
