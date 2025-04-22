package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

const PRELUDE_TYPENAME_VECINT = "prelude.vector.int"
const PRELUDE_TYPENAME_VECSTR = "prelude.vector.string"
const PRELUDE_TYPENAME_VECBOOL = "prelude.vector.bool"

func declareVectorType(ir *ir.Module, Prelude *mTypes.PreludeProps) {
	// i32 type vector
	// 型: struct { i32* %arrElm, i64 %len}
	vectorIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	vectorIntType.SetName(PRELUDE_TYPENAME_VECINT)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECINT, vectorIntType)
	Prelude.Types.VectorInt = vectorIntType

	// i1 (bool) type vector
	// 型: struct { i1* %arrElm, i64 %len}
	vectorBoolType := types.NewStruct(types.NewPointer(types.I1), types.I64)
	vectorBoolType.SetName(PRELUDE_TYPENAME_VECBOOL)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECBOOL, vectorBoolType)
	Prelude.Types.VectorBool = vectorBoolType

	// i8* (string) type vector
	// 型: struct { i8** %arrElm, i64 %len}
	vectorStringType := types.NewStruct(types.NewPointer(types.I8Ptr), types.I64)
	vectorStringType.SetName(PRELUDE_TYPENAME_VECSTR)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECSTR, vectorStringType)
	Prelude.Types.VectorString = vectorStringType

}

// TODO: stop to use
func GetLLVMTypeFromString(typeName string, prelude *mTypes.PreludeProps) (types.Type, types.Type) {

	switch typeName {

	case PRELUDE_TYPENAME_VECINT:
		return prelude.Types.VectorInt, types.I32

	case PRELUDE_TYPENAME_VECSTR:
		return prelude.Types.VectorString, types.I8Ptr

	case PRELUDE_TYPENAME_VECBOOL:
		return prelude.Types.VectorBool, types.I1
	}

	return nil, nil

}

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
	ctx.block.NewStore(constant.NewInt(types.I64, 0), loopIndex)

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

	incI := loopBlock.NewAdd(idx, constant.NewInt(types.I64, 1))
	loopBlock.NewStore(incI, loopIndex)
	loopBlock.NewBr(condBlock)

	ctx.block = endBlock

	return newVecPtr

}
