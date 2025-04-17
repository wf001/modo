package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludePop(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	var oldArrType *types.ArrayType

	var oldArrPtr value.Value

	switch v := node.IRValue.(type) {
	case *ir.InstCall, *ir.InstBitCast:
		oldArrType, _ = mTypes.AssertArrType(v)
		oldArrPtr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", node.IRValue)
	}

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(oldArrType.ElemType)))
	newArrSize := constant.NewInt(types.I64, int64(oldArrType.Len-1))
	allocSize := block.NewMul(typeSize, newArrSize)
	newArrType := types.NewArray(oldArrType.Len-1, oldArrType.ElemType)
	newArrPtr := CopyArray(block, internal, allocSize, oldArrType, oldArrPtr, newArrType)

	return newArrPtr
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func PreludePopOld(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {

	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len-1, oldArr.Len-1)

	return newArr
}
