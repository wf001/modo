package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeGet(ctx *Context, n *mTypes.Node) value.Value {
	tyPtr, isTyPtr := n.IRValue.Type().(*types.PointerType)
	tyStr, _ := tyPtr.ElemType.(*types.StructType)

	specifiedFields := ctx.prog.Declare.Type.Struct[tyStr.TypeName].Field[n.Next.Val]
	i := specifiedFields.Pos
	structeType, elemType := ctx.prog.Declare.Type.Struct[tyStr.TypeName].Types, specifiedFields.Type

	targetStructPtr := n.IRValue

	if structeType == nil || elemType == nil {
		return constant.NewNull(types.NewPointer(types.I32))
	}

	if !isTyPtr {
		loadedStruct := ctx.block.NewLoad(structeType, targetStructPtr)
		elemPtr := ctx.block.NewExtractValue(loadedStruct, i)
		return elemPtr
	}

	cond := ctx.block.NewICmp(
		enum.IPredEQ,
		targetStructPtr,
		constant.NewNull(tyPtr),
	)

	nullBlock := ctx.NewBlock("get.null", n)
	nonNullBlock := ctx.NewBlock("get.non.null", n)
	mergeBlock := ctx.NewBlock("get.merge", n)

	ctx.block.NewCondBr(cond, nullBlock, nonNullBlock)

	nullBlock.NewBr(mergeBlock)

	loadedStruct := nonNullBlock.NewLoad(structeType, targetStructPtr)
	elemPtr := nonNullBlock.NewExtractValue(loadedStruct, i)
	nonNullBlock.NewBr(mergeBlock)

	ctx.block = mergeBlock

	t, isElemPtr := elemPtr.Type().(*types.PointerType)
	if !isElemPtr {
		loadedStruct := ctx.block.NewLoad(structeType, targetStructPtr)
		elemPtr := ctx.block.NewExtractValue(loadedStruct, i)
		return elemPtr
	}

	result := mergeBlock.NewPhi(
		[]*ir.Incoming{
			{X: constant.NewNull(t), Pred: nullBlock},
			{X: elemPtr, Pred: nonNullBlock},
		}...,
	)

	ctx.block = mergeBlock
	return result
}
