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

func InvokeAssoc(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var oldArrType *types.ArrayType

	var oldArrPtr value.Value

	switch v := node.IRValue.(type) {
	case *ir.InstCall, *ir.InstBitCast:
		oldArrType = mTypes.GetArrType(v)
		oldArrPtr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", node.IRValue)
	}

	pos, newValue := node.Next.IRValue, node.Next.Next.IRValue

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(oldArrType.ElemType)))
	newArrSize := constant.NewInt(types.I64, int64(oldArrType.Len))
	allocSize := block.NewMul(typeSize, newArrSize)
	newArrType := types.NewArray(oldArrType.Len, oldArrType.ElemType)
	newArrPtr := vector.CopyArray(block, libs, allocSize, oldArrType, oldArrPtr, newArrType)

	destPtr := block.NewGetElementPtr(
		newArrType,
		newArrPtr,
		constant.NewInt(types.I32, 0),
		pos,
	)
	block.NewStore(newValue, destPtr)

	return newArrPtr
}
func InvokeAssocOld(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	pos, newValue := node.Next.IRValue, node.Next.Next.IRValue
	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := vector.CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len, oldArr.Len)

	newElemPtr := block.NewGetElementPtr(
		newArr.ElemType,
		newArr,
		constant.NewInt(types.I32, 0),
		pos,
	)
	block.NewStore(newValue, newElemPtr)

	return newArr
}
