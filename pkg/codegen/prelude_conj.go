package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeConj(ctx *Context, n *mTypes.Node) value.Value {
	// もとのベクターの構造体をロード
	oldStructedVecPtr := n.IRValue
	structPtrType := oldStructedVecPtr.Type().(*types.PointerType)
	structType := structPtrType.ElemType.(*types.StructType)

	// array type の情報を取得
	structedVecType, elemType := GetLLVMTypeFromString(structType.TypeName, ctx.prog.Prelude)

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)

	// 長さ +1 の新しい vector を確保
	elemSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newLen := ctx.block.NewAdd(oldLen, constant.NewInt(types.I64, 1))
	newVecAllocSize := ctx.block.NewMul(elemSize, newLen)

	newVecAllocPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, newVecAllocSize)
	newVecPtr := ctx.block.NewBitCast(newVecAllocPtr, types.NewPointer(elemType))

	// もとの要素をコピー
	loopIndex := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(constant.NewInt(types.I64, 0), loopIndex)

	loopBlock := ctx.function.NewBlock(n.GetBlockName("copy.loop", ctx.function.Blocks))
	condBlock := ctx.function.NewBlock(n.GetBlockName("copy.cond", ctx.function.Blocks))
	endBlock := ctx.function.NewBlock(n.GetBlockName("copy.end", ctx.function.Blocks))

	ctx.block.NewBr(condBlock)

	// 条件チェックブロック
	idx := condBlock.NewLoad(types.I64, loopIndex)
	copyContinue := condBlock.NewICmp(enum.IPredULT, idx, oldLen)
	condBlock.NewCondBr(copyContinue, loopBlock, endBlock)

	// コピー処理ブロック
	oldArrElemPtr := loopBlock.NewGetElementPtr(elemType, oldVecPtr, idx)
	oldElem := loopBlock.NewLoad(elemType, oldArrElemPtr)
	newVecElemPtr := loopBlock.NewGetElementPtr(elemType, newVecPtr, idx)
	loopBlock.NewStore(oldElem, newVecElemPtr)

	incI := loopBlock.NewAdd(idx, constant.NewInt(types.I64, 1))
	loopBlock.NewStore(incI, loopIndex)
	loopBlock.NewBr(condBlock)

	// end: 追加要素を挿入
	ctx.block = endBlock
	newVecElemPtr = ctx.block.NewGetElementPtr(elemType, newVecPtr, oldLen)
	ctx.block.NewStore(n.Next.IRValue, newVecElemPtr)

	// 新しい構造体を返す
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
