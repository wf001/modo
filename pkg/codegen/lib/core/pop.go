package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib/core/vector"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokePop(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {

	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := vector.CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len-1, oldArr.Len-1)

	return newArr
}
