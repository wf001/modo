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

func PreludeReduce(ctx *Context, n *mTypes.Node) value.Value {
	if !n.IsKind(mTypes.ND_VAR_REFERENCE) {
		log.Panic("%s: unexpected character used", error.ERROR_SYNTAX_ERROR)
	}
	var f *ir.Func

	// find declared function
	// TODO: find also in prelude function
	for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
		if ctx.mod.Funcs[i].GlobalName == n.GetFuncName() {
			f = ctx.mod.Funcs[i]
			break
		}
	}

	initialValue := n.Next.IRValue
	oldStructedVecPtr := n.Next.Next.IRValue

	structedVecType, _ := mTypes.GetVectorTypeFromPtr(oldStructedVecPtr)
	elemType := f.Sig.RetType

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)

	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)

	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))

	newAllocPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, elemSize)
	newPtr := ctx.block.NewBitCast(newAllocPtr, types.NewPointer(elemType))

	loopIndexPtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, loopIndexPtr)

	condBlock := ctx.NewBlock("reduce.cond", n)
	loopBlock := ctx.NewBlock("reduce.loop", n)
	endBlock := ctx.NewBlock("reduce.end", n)

	// f(initialValue, v[0])
	loopIdx := ctx.block.NewLoad(types.I64, loopIndexPtr)
	oldArrElemPtr := ctx.block.NewGetElementPtr(elemType, oldVecPtr, loopIdx)
	oldElem := ctx.block.NewLoad(elemType, oldArrElemPtr)
	result1 := ctx.block.NewCall(f, initialValue, oldElem)
	ctx.block.NewStore(result1, newPtr)
	ctx.block.NewStore(
		ctx.block.NewAdd(loopIdx, mTypes.I64one),
		loopIndexPtr,
	)
	ctx.block.NewBr(condBlock)

	// loop start
	loopIdx = condBlock.NewLoad(types.I64, loopIndexPtr)
	isIdxLTOldLen := condBlock.NewICmp(enum.IPredULT, loopIdx, oldLen)
	condBlock.NewCondBr(isIdxLTOldLen, loopBlock, endBlock)

	oldArrElemPtr = loopBlock.NewGetElementPtr(elemType, oldVecPtr, loopIdx)
	oldElem = loopBlock.NewLoad(elemType, oldArrElemPtr)
	prevResult := loopBlock.NewLoad(elemType, newPtr)
	result1 = loopBlock.NewCall(f, prevResult, oldElem)
	loopBlock.NewStore(result1, newPtr)

	loopBlock.NewStore(
		loopBlock.NewAdd(loopIdx, mTypes.I64one),
		loopIndexPtr,
	)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock
	// loop end

	return newPtr
}
