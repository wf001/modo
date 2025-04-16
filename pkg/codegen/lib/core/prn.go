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
	rootTy := n.IRValue.Type()
	formatStr, _ := mTypes.GetPrintFormat(rootTy, libs)

	if rootTy.Equal(types.I1) {
		value = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)

	} else if rootTy.Equal(types.Void) {
		value = libs.GlobalVar.NilValue

	} else if pointerElemTy, isPtr := rootTy.(*types.PointerType); isPtr {
		if isArr := pointerElemTy.ElemType.Equal(types.I32); isArr {
			ptr := block.NewLoad(pointerElemTy, n.IRValue)
			value = block.NewGetElementPtr(
				pointerElemTy,
				ptr,
				constant.NewInt(types.I32, 0),
			)
			formatStr, _ = mTypes.GetPrintFormat(pointerElemTy.ElemType, libs)
		}
	}

	block.NewCall(libs.Printf.FuncPtr, formatStr, value)
}

func prnVector(libs *mTypes.BuiltinLibProp, block *ir.Block, n *mTypes.Node) {
	var arr value.Value
	var arrType *types.ArrayType

	switch v := n.IRValue.(type) {
	// The vector on global space: InstCall,
	// The vector on local space: InstBitCast
	case *ir.InstCall, *ir.InstBitCast:
		arrType = mTypes.GetArrType(v)
		arr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", n.IRValue)
	}

	block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatBracketOpen)

	for i := uint64(0); i < arrType.Len; i++ {
		elem := mTypes.LoadArrElem(block, arr, arrType, i)

		formatStr, _ := mTypes.GetPrintFormat(elem.ElemType, libs)
		if arrType.ElemType == types.I1 {
			v := block.NewSelect(elem, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)
			block.NewCall(libs.Printf.FuncPtr, formatStr, v)
		} else {
			block.NewCall(libs.Printf.FuncPtr, formatStr, elem)
		}

		if i < uint64(arrType.Len-1) {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatComma)
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}
	}

	block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatBracketClose)
}

func InvokePrn(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {

	for n := node; n != nil; n = n.Next {
		rootTy := n.IRValue.Type()

		if mTypes.IsScalar(n.IRValue) {
			prnScalar(libs, block, n)

		} else if pointerElemTy, ok := rootTy.(*types.PointerType); ok {
			if _, ok := pointerElemTy.ElemType.(*types.ArrayType); ok {
				prnVector(libs, block, n)
			} else {
				prnScalar(libs, block, n)
			}
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
