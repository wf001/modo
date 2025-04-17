package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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
	rootTy := value.Type()
	formatStr, _ := mTypes.GetPrintFormat(rootTy, libs)

	if rootTy.Equal(types.I1) {
		value = block.NewSelect(value, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)

	} else if rootTy.Equal(types.Void) {
		value = libs.GlobalVar.NilValue

	} else if pointerElemTy, isPtr := rootTy.(*types.PointerType); isPtr {
		if isStr := rootTy.Equal(types.I8Ptr); !isStr {
			ptr := block.NewLoad(pointerElemTy, value)
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
		arrType, _ = mTypes.AssertArrType(v)
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

func prnStructVector(
	fnc *ir.Func,
	libs *mTypes.BuiltinLibProp,
	block *ir.Block,
	n *mTypes.Node,
	arrType *mTypes.ArrayTypeProps,
) {
	loaded := block.NewLoad(arrType.TypeInt, n.IRValue)

	// 構造体のフィールドから arrPtr と len を取り出す
	resultArrPtr := block.NewExtractValue(loaded, 0)
	resultLen := block.NewExtractValue(loaded, 1)

	// インデックスの初期化
	idx := block.NewAlloca(types.I64)
	idx.SetName("idx")
	block.NewStore(constant.NewInt(types.I64, 0), idx)

	loopBlock := fnc.NewBlock("loop")
	endBlock := fnc.NewBlock("end")
	block.NewBr(loopBlock)

	// ループ内部の処理
	// i をロード
	i := loopBlock.NewLoad(types.I64, idx)

	// 配列の各要素を取り出して表示
	elemPtr := loopBlock.NewGetElementPtr(types.I32, resultArrPtr, i)
	elem := loopBlock.NewLoad(types.I32, elemPtr)
	loopBlock.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatDigit, elem)

	// i++
	nextI := loopBlock.NewAdd(i, constant.NewInt(types.I64, 1))
	loopBlock.NewStore(nextI, idx)

	// i < len ?
	cond := loopBlock.NewICmp(enum.IPredSLT, nextI, resultLen)
	loopBlock.NewCondBr(cond, loopBlock, endBlock)
	block = endBlock

}

func InvokePrn(
	fnc *ir.Func,
	block *ir.Block,
	libs *mTypes.BuiltinLibProp,
	node *mTypes.Node,
	arrType *mTypes.ArrayTypeProps,
) value.Value {

	for n := node; n != nil; n = n.Next {
		if mTypes.IsScalar(n.IRValue) {
			prnScalar(libs, block, n)

		} else if _, ok := n.IRValue.Type().(*types.PointerType); ok {
			prnStructVector(fnc, libs, block, n, arrType)

		} else if _, ok := mTypes.AssertArrType(n.IRValue); ok {
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
