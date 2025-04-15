package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib/core/vector"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeConj(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var arrType *types.ArrayType

	var oldArrPtr value.Value

	switch v := node.IRValue.(type) {
	case *ir.InstCall, *ir.InstBitCast:
		arrType = mTypes.GetArrType(v)
		oldArrPtr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", node.IRValue)
	}

	newValue := node.Next.IRValue

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(arrType.ElemType)))
	newArrSize := constant.NewInt(types.I64, int64(arrType.Len+1))
	allocSize := block.NewMul(typeSize, newArrSize)
	allocatedPtr := block.NewCall(libs.Malloc.FuncPtr, allocSize)

	newArrType := types.NewArray(arrType.Len+1, arrType.ElemType)
	newArrPtr := block.NewBitCast(allocatedPtr, types.NewPointer(newArrType))

	for i := uint64(0); i < arrType.Len; i++ {
		elem := mTypes.LoadArrElem(block, oldArrPtr, arrType, i)
		destPtr := block.NewGetElementPtr(
			newArrType,
			newArrPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		block.NewStore(elem, destPtr)
	}

	destPtr := block.NewGetElementPtr(
		newArrType,
		newArrPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(arrType.Len)),
	)
	block.NewStore(newValue, destPtr)

	return newArrPtr
}
func InvokeConjOld(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	newValue := node.Next.IRValue
	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := vector.CopyArray(block, oldArrPtr, oldArr.ElemType, oldArr.Len+1, oldArr.Len)

	newElemPtr := block.NewGetElementPtr(
		newArr.ElemType,
		newArr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(oldArr.Len)),
	)
	block.NewStore(node.Next.IRValue, newElemPtr)
	block.NewStore(newValue, newElemPtr)

	return newArr
}
