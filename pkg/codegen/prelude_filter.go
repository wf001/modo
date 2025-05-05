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

func PreludeFilter(ctx *Context, n *mTypes.Node) value.Value {
	if !n.IsKind(mTypes.ND_VAR_REFERENCE) {
		log.Panic("%s: unexpected character used", error.ERROR_SYNTAX_ERROR)
	}
	var pred *ir.Func

	// find declared function
	// TODO: find also in prelude function
	for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
		if ctx.mod.Funcs[i].GlobalName == n.GetFuncName() {
			pred = ctx.mod.Funcs[i]
			break
		}
	}

	oldStructedVecPtr := n.Next.IRValue
	structedVecType := mTypes.GetStructTypeFromPtr(oldStructedVecPtr)
	elemType := structedVecType.Fields[0].(*types.PointerType).ElemType

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)

	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)

	// get the number of elements for which the predicate returns true
	// This step is necessary to allocate the correct number of bytes with malloc
	numTruePtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, numTruePtr)
	countLoopIdxPtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, countLoopIdxPtr)

	countCondBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.count.cond", ctx.function.Blocks),
	)
	countLoopBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.count.loop", ctx.function.Blocks),
	)
	countLoopIncCountBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.count.loop.inc.count", ctx.function.Blocks),
	)
	countLoopIncIdxBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.count.loop.inc.idx", ctx.function.Blocks),
	)
	countEndBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.count.end", ctx.function.Blocks),
	)

	ctx.block.NewBr(countCondBlock)

	numTrue := countCondBlock.NewLoad(types.I64, numTruePtr)
	numTrue.SetName(n.GetVarName("filter.count.true", ctx.block.Insts))
	countloopIdx := countCondBlock.NewLoad(types.I64, countLoopIdxPtr)
	countloopIdx.SetName(n.GetVarName("filter.count.loop.idx", ctx.block.Insts))

	isIdxLTOldLen := countCondBlock.NewICmp(enum.IPredULT, countloopIdx, oldLen)
	countCondBlock.NewCondBr(isIdxLTOldLen, countLoopBlock, countEndBlock)

	oldElemPtr := countLoopBlock.NewGetElementPtr(elemType, oldVecPtr, countloopIdx)
	oldElem := countLoopBlock.NewLoad(elemType, oldElemPtr)
	isPredResultTrue := countLoopBlock.NewCall(pred, oldElem)
	countLoopBlock.NewCondBr(isPredResultTrue, countLoopIncCountBlock, countLoopIncIdxBlock)

	countLoopIncCountBlock.NewStore(
		countLoopIncCountBlock.NewAdd(numTrue, mTypes.I64one),
		numTruePtr,
	)
	countLoopIncCountBlock.NewBr(countLoopIncIdxBlock)

	countLoopIncIdxBlock.NewStore(
		countLoopIncIdxBlock.NewAdd(countloopIdx, mTypes.I64one),
		countLoopIdxPtr,
	)
	countLoopIncIdxBlock.NewBr(countCondBlock)

	ctx.block = countEndBlock

	// Create a vector containing only elements for which the predicate returns true
	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newVecAllocSize := ctx.block.NewMul(elemSize, numTrue)

	newVecAllocPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, newVecAllocSize)
	newVecPtr := ctx.block.NewBitCast(newVecAllocPtr, types.NewPointer(elemType))

	conjLoopIdxPtr := ctx.block.NewAlloca(types.I64)
	newVecIdxPtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, conjLoopIdxPtr)
	ctx.block.NewStore(mTypes.I64zero, newVecIdxPtr)

	condBlock := ctx.function.NewBlock(n.GetBlockName("filter.new.vec.cond", ctx.function.Blocks))
	loopBlock := ctx.function.NewBlock(n.GetBlockName("filter.new.vec.loop", ctx.function.Blocks))
	loopConjBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.new.vec.loop.conj", ctx.function.Blocks),
	)
	loopIncBlock := ctx.function.NewBlock(
		n.GetBlockName("filter.new.vec.loop.inc", ctx.function.Blocks),
	)
	endBlock := ctx.function.NewBlock(n.GetBlockName("filter.new.vec.end", ctx.function.Blocks))

	ctx.block.NewBr(condBlock)

	conjLoopIdx := condBlock.NewLoad(types.I64, conjLoopIdxPtr)
	newVecIdx := condBlock.NewLoad(types.I64, newVecIdxPtr)

	isIdxLTOldLen = condBlock.NewICmp(enum.IPredULT, conjLoopIdx, oldLen)
	condBlock.NewCondBr(isIdxLTOldLen, loopBlock, endBlock)

	// Note: redeclare needed to avoid 'Instruction does not dominate all uses!' error
	oldElemPtr = loopBlock.NewGetElementPtr(elemType, oldVecPtr, conjLoopIdx)
	oldElem = loopBlock.NewLoad(elemType, oldElemPtr)
	isPredResultTrue = loopBlock.NewCall(pred, oldElem)
	loopBlock.NewCondBr(isPredResultTrue, loopConjBlock, loopIncBlock)

	newVecElemPtr := loopConjBlock.NewGetElementPtr(elemType, newVecPtr, newVecIdx)
	loopConjBlock.NewStore(oldElem, newVecElemPtr)
	loopConjBlock.NewStore(
		loopConjBlock.NewAdd(newVecIdx, mTypes.I64one),
		newVecIdxPtr,
	)
	loopConjBlock.NewBr(loopIncBlock)

	loopIncBlock.NewStore(
		loopIncBlock.NewAdd(conjLoopIdx, mTypes.I64one),
		conjLoopIdxPtr,
	)
	loopIncBlock.NewBr(condBlock)

	ctx.block = endBlock
	// loop end

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
	ctx.block.NewStore(numTrue, newLenField)

	return newStructAlloca
}
