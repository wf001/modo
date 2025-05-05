package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func CopyVector(
	ctx *Context,
	oldStructedVec *ir.InstLoad,
	n *mTypes.Node,
	elemType types.Type,
	oldLen *ir.InstExtractValue,
	newLen value.Value,
) *ir.InstBitCast {
	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)

	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newVecAllocSize := ctx.block.NewMul(elemSize, newLen)

	newVecAllocPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, newVecAllocSize)
	newVecPtr := ctx.block.NewBitCast(newVecAllocPtr, types.NewPointer(elemType))

	loopIndexPtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, loopIndexPtr)

	loopBlock := ctx.NewBlock("copy.loop", n)
	condBlock := ctx.NewBlock("copy.cond", n)
	endBlock := ctx.NewBlock("copy.end", n)

	ctx.block.NewBr(condBlock)

	loopIdx := condBlock.NewLoad(types.I64, loopIndexPtr)
	isIdxLTOldLen := condBlock.NewICmp(enum.IPredULT, loopIdx, oldLen)
	condBlock.NewCondBr(isIdxLTOldLen, loopBlock, endBlock)

	oldArrElemPtr := loopBlock.NewGetElementPtr(elemType, oldVecPtr, loopIdx)
	oldElem := loopBlock.NewLoad(elemType, oldArrElemPtr)
	newVecElemPtr := loopBlock.NewGetElementPtr(elemType, newVecPtr, loopIdx)
	loopBlock.NewStore(oldElem, newVecElemPtr)

	loopBlock.NewStore(
		loopBlock.NewAdd(loopIdx, mTypes.I64one),
		loopIndexPtr,
	)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock

	return newVecPtr

}
