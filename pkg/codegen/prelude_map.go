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

	// find declared function
	// TODO: find also in prelude function
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

	// create a new vector which element are applied by specified function
	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)

	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newVecAllocSize := ctx.block.NewMul(elemSize, oldLen)

	newVecAllocPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, newVecAllocSize)
	newVecPtr := ctx.block.NewBitCast(newVecAllocPtr, types.NewPointer(elemType))

	loopIndexPtr := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(mTypes.I64zero, loopIndexPtr)

	loopBlock := ctx.function.NewBlock(n.GetBlockName("map.loop", ctx.function.Blocks))
	condBlock := ctx.function.NewBlock(n.GetBlockName("map.cond", ctx.function.Blocks))
	endBlock := ctx.function.NewBlock(n.GetBlockName("map.end", ctx.function.Blocks))

	ctx.block.NewBr(condBlock)

	loopIdx := condBlock.NewLoad(types.I64, loopIndexPtr)
	isIdxLTOldLen := condBlock.NewICmp(enum.IPredULT, loopIdx, oldLen)
	condBlock.NewCondBr(isIdxLTOldLen, loopBlock, endBlock)

	oldArrElemPtr := loopBlock.NewGetElementPtr(elemType, oldVecPtr, loopIdx)
	oldElem := loopBlock.NewLoad(elemType, oldArrElemPtr)
	newElem := loopBlock.NewCall(f, oldElem)
	newVecElemPtr := loopBlock.NewGetElementPtr(elemType, newVecPtr, loopIdx)
	loopBlock.NewStore(newElem, newVecElemPtr)

	loopBlock.NewStore(
		loopBlock.NewAdd(loopIdx, mTypes.I64one),
		loopIndexPtr,
	)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock
	// loop end

	newStructVecType := mTypes.DeclareVectorType(
		ctx.mod,
		elemType,
		&mTypes.NodeType{Value: mTypes.TY_VECTOR},
	)
	newStructAlloca := ctx.block.NewAlloca(newStructVecType)

	newVecField := ctx.block.NewGetElementPtr(
		newStructVecType,
		newStructAlloca,
		mTypes.I32zero,
		mTypes.I32zero,
	)
	ctx.block.NewStore(newVecPtr, newVecField)

	newLenField := ctx.block.NewGetElementPtr(
		newStructVecType,
		newStructAlloca,
		mTypes.I32zero,
		mTypes.I32one,
	)
	ctx.block.NewStore(oldLen, newLenField)

	return newStructAlloca
}
