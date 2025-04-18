package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeConj(ctx *Context, n *mTypes.Node) value.Value {
	oldStructedVecPtr := n.IRValue
	structPtrType := oldStructedVecPtr.Type().(*types.PointerType)
	structType := structPtrType.ElemType.(*types.StructType)

	structedVecType, elemType := GetLLVMTypeFromString(structType.TypeName, ctx.prog.Prelude)

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

// Note: remain here until it will be defined the strategy of memory lifecycle
func PreludeConjOld(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	oldArrPtr, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", oldArrPtr)
	}

	newValue := node.Next.IRValue
	oldArr := oldArrPtr.ElemType.(*types.ArrayType)

	newArr := CopyArrayOld(block, oldArrPtr, oldArr.ElemType, oldArr.Len+1, oldArr.Len)

	newElemPtr := block.NewGetElementPtr(
		newArr.ElemType,
		newArr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(oldArr.Len)),
	)
	block.NewStore(node.Next.IRValue, newElemPtr)
	block.NewStore(newValue, newElemPtr)

	return newArr
}
