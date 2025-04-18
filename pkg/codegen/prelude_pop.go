package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludePop(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	return nil
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
