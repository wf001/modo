package codegen

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeConj(ctx *Context, n *mTypes.Node) value.Value {
	oldStructedVecPtr := n.IRValue

	structedVecType := mTypes.GetStructTypeFromPtr(oldStructedVecPtr)
	elemType := structedVecType.Fields[0].(*types.PointerType).ElemType

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)
	newLen := ctx.block.NewAdd(oldLen, constant.NewInt(types.I64, 1))

	newVecPtr := CopyVector(ctx, oldStructedVec, n, elemType, oldLen, newLen)

	newVecElemPtr := ctx.block.NewGetElementPtr(elemType, newVecPtr, oldLen)
	newVal := n.Next.IRValue
	ctx.block.NewStore(newVal, newVecElemPtr)

	newStructAlloca := ctx.block.NewAlloca(structedVecType)

	newVecField := ctx.block.NewGetElementPtr(
		structedVecType,
		newStructAlloca,
		newI32("0"),
		newI32("0"),
	)
	ctx.block.NewStore(newVecPtr, newVecField)

	newLenField := ctx.block.NewGetElementPtr(
		structedVecType,
		newStructAlloca,
		newI32("0"),
		newI32("1"),
	)
	ctx.block.NewStore(newLen, newLenField)

	return newStructAlloca
}
