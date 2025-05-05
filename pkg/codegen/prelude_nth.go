package codegen

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeNth(ctx *Context, n *mTypes.Node) value.Value {
	i, _ := strconv.ParseInt(n.Next.Val, 10, 32)
	idx := constant.NewInt(types.I64, i)

	structedVecPtr := n.IRValue
	structedVecType, elemType := mTypes.GetVectorTypeFromPtr(structedVecPtr)

	nullPtr := constant.NewNull(types.NewPointer(elemType))

	if i < 0 {
		return nullPtr
	}

	loadedStructedVec := ctx.block.NewLoad(structedVecType, structedVecPtr)
	loadedVecPtr := ctx.block.NewExtractValue(loadedStructedVec, 0)
	loadedLength := ctx.block.NewExtractValue(loadedStructedVec, 1)
	maxIdx := ctx.block.NewSub(loadedLength, mTypes.I64one)

	isIdxOutOfRange := ctx.block.NewICmp(enum.IPredSGT, idx, maxIdx)
	isIdxOutOfRange.SetName(n.GetVarName("get.is.idx.out.of.range"))

	inRangeBlock := ctx.NewBlock("get.idx.in.range", n)
	outOfRangeBlock := ctx.NewBlock("get.idx.out.of.range", n)
	mergeBlock := ctx.NewBlock("get.idx.merge", n)

	ctx.block.NewCondBr(isIdxOutOfRange, outOfRangeBlock, inRangeBlock)

	// branched when specified index is in range
	ctx.block = inRangeBlock
	vecElemPtr := ctx.block.NewGetElementPtr(elemType, loadedVecPtr, idx)
	elem := ctx.block.NewLoad(elemType, vecElemPtr)
	okPtr := ctx.block.NewAlloca(elemType)
	ctx.block.NewStore(elem, okPtr)
	ctx.block.NewBr(mergeBlock)

	// branched when specified index is out of range
	ctx.block = outOfRangeBlock
	ctx.block.NewBr(mergeBlock)

	ctx.block = mergeBlock
	incs := []*ir.Incoming{
		ir.NewIncoming(okPtr, inRangeBlock),
		ir.NewIncoming(nullPtr, outOfRangeBlock),
	}
	result := ctx.block.NewPhi(incs...)

	return result
}
