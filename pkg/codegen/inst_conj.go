package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeConj(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var oldArrType *types.ArrayType

	var oldArrPtr value.Value

	switch v := node.IRValue.(type) {
	case *ir.InstCall, *ir.InstBitCast, *ir.Param:
		oldArrType, _ = mTypes.AssertArrType(v)
		oldArrPtr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", node.IRValue)
	}

	newValue := node.Next.IRValue

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(oldArrType.ElemType)))
	newArrSize := constant.NewInt(types.I64, int64(oldArrType.Len+1))
	allocSize := block.NewMul(typeSize, newArrSize)
	newArrType := types.NewArray(oldArrType.Len+1, oldArrType.ElemType)
	newArrPtr := CopyArray(block, libs, allocSize, oldArrType, oldArrPtr, newArrType)

	destPtr := block.NewGetElementPtr(
		newArrType,
		newArrPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(oldArrType.Len)),
	)
	block.NewStore(newValue, destPtr)

	return newArrPtr
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func InvokeConjOld(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	newValue := node.Next.IRValue
	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len+1, oldArr.Len)

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
