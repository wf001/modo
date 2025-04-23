package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeAssoc(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	return nil
}
func PreludeAssocOld(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	pos, newValue := node.Next.IRValue, node.Next.Next.IRValue
	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len, oldArr.Len)

	newElemPtr := block.NewGetElementPtr(
		newArr.ElemType,
		newArr,
		mTypes.I32zero,
		pos,
	)
	block.NewStore(newValue, newElemPtr)

	return newArr
}
