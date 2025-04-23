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

	loopIndex := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, loopIndex)

	loopBlock := ctx.function.NewBlock(n.GetBlockName("copy.loop", ctx.function.Blocks))
	condBlock := ctx.function.NewBlock(n.GetBlockName("copy.cond", ctx.function.Blocks))
	endBlock := ctx.function.NewBlock(n.GetBlockName("copy.end", ctx.function.Blocks))

	ctx.block.NewBr(condBlock)

	idx := condBlock.NewLoad(types.I64, loopIndex)
	copyContinue := condBlock.NewICmp(enum.IPredULT, idx, oldLen)
	condBlock.NewCondBr(copyContinue, loopBlock, endBlock)

	oldArrElemPtr := loopBlock.NewGetElementPtr(elemType, oldVecPtr, idx)
	oldElem := loopBlock.NewLoad(elemType, oldArrElemPtr)
	newVecElemPtr := loopBlock.NewGetElementPtr(elemType, newVecPtr, idx)
	loopBlock.NewStore(oldElem, newVecElemPtr)

	incI := loopBlock.NewAdd(idx, mTypes.I64one)
	loopBlock.NewStore(incI, loopIndex)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock

	return newVecPtr

}
