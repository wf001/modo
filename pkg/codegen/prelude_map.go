package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeMap(ctx *Context, n *mTypes.Node) value.Value {
	if !n.IsKind(mTypes.ND_VAR_REFERENCE) {
		log.Panic("%s: unexpected character used", error.ERROR_SYNTAX_ERROR)
	}
	var f *ir.Func

	for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
		if ctx.mod.Funcs[i].GlobalName == n.GetFuncName() {
			f = ctx.mod.Funcs[i]
			break
		}
	}

	oldStructedVecPtr := n.Next.IRValue
	structedVecType := mTypes.GetStructTypeFromPtr(oldStructedVecPtr)
	elemType := f.Sig.RetType

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)

	/////
	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)

	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newVecAllocSize := ctx.block.NewMul(elemSize, oldLen)

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
	newElem := loopBlock.NewCall(f, oldElem)
	newVecElemPtr := loopBlock.NewGetElementPtr(elemType, newVecPtr, idx)
	loopBlock.NewStore(newElem, newVecElemPtr)

	incI := loopBlock.NewAdd(idx, mTypes.I64one)
	loopBlock.NewStore(incI, loopIndex)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock
	/////

	newStructAlloca := ctx.block.NewAlloca(structedVecType)

	newVecField := ctx.block.NewGetElementPtr(
		structedVecType,
		newStructAlloca,
		mTypes.I32zero,
		mTypes.I32zero,
	)
	ctx.block.NewStore(newVecPtr, newVecField)

	newLenField := ctx.block.NewGetElementPtr(
		structedVecType,
		newStructAlloca,
		mTypes.I32zero,
		mTypes.I32one,
	)
	ctx.block.NewStore(oldLen, newLenField)

	return newStructAlloca
}
