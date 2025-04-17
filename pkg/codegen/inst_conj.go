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

func InvokeConj(ctx *Context, n *mTypes.Node) value.Value {
	// もとのベクターの構造体をロード
	originalVector := n.IRValue
	structPtrType := originalVector.Type().(*types.PointerType)
	structType := structPtrType.ElemType.(*types.StructType)

	var m1 = map[string]*types.StructType{
		"array.int":    ctx.prog.ArrayType.TypeInt,
		"array.string": ctx.prog.ArrayType.TypeString,
	}
	var m2 = map[string]types.Type{
		"array.int":    types.I32,
		"array.string": types.I8Ptr,
	}
	// array type の情報を取得
	structedArrType := m1[structType.TypeName]
	elemType := m2[structType.TypeName]

	loaded := ctx.block.NewLoad(structedArrType, originalVector)
	resultArrPtr := ctx.block.NewExtractValue(loaded, 0)
	resultLen := ctx.block.NewExtractValue(loaded, 1)

	// 長さ +1 の新しい vector を確保
	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	newLen := ctx.block.NewAdd(resultLen, constant.NewInt(types.I64, 1))
	allocSize := ctx.block.NewMul(typeSize, newLen)

	allocatedPtr := ctx.block.NewCall(ctx.prog.BuiltinLibs.Malloc.FuncPtr, allocSize)
	newArrayPtr := ctx.block.NewBitCast(allocatedPtr, types.NewPointer(elemType))

	// もとの要素をコピー
	loopIndex := ctx.block.NewAlloca(types.I64)
	ctx.block.NewStore(constant.NewInt(types.I64, 0), loopIndex)

	loop := ctx.function.NewBlock(n.GetBlockName("copy_loop", ctx.function.Blocks))
	cond := ctx.function.NewBlock(n.GetBlockName("copy_cond", ctx.function.Blocks))
	end := ctx.function.NewBlock(n.GetBlockName("copy_end", ctx.function.Blocks))

	ctx.block.NewBr(cond)

	// 条件チェックブロック
	condI := cond.NewLoad(types.I64, loopIndex)
	shouldContinue := cond.NewICmp(enum.IPredULT, condI, resultLen)
	cond.NewCondBr(shouldContinue, loop, end)

	// コピー処理ブロック
	elemPtr := loop.NewGetElementPtr(elemType, resultArrPtr, condI)
	elem := loop.NewLoad(elemType, elemPtr)
	newElemPtr := loop.NewGetElementPtr(elemType, newArrayPtr, condI)
	loop.NewStore(elem, newElemPtr)

	nextI := loop.NewAdd(condI, constant.NewInt(types.I64, 1))
	loop.NewStore(nextI, loopIndex)
	loop.NewBr(cond)

	// end: 追加要素を挿入
	ctx.block = end
	newElemPtr = ctx.block.NewGetElementPtr(elemType, newArrayPtr, resultLen)
	ctx.block.NewStore(n.Next.IRValue, newElemPtr)

	// 新しい構造体を返す
	newStructAlloca := ctx.block.NewAlloca(structedArrType)

	ptrField := ctx.block.NewGetElementPtr(
		structedArrType,
		newStructAlloca,
		newI32("0"),
		newI32("0"),
	)
	ctx.block.NewStore(newArrayPtr, ptrField)

	lenField := ctx.block.NewGetElementPtr(
		structedArrType,
		newStructAlloca,
		newI32("0"),
		newI32("1"),
	)
	ctx.block.NewStore(newLen, lenField)

	return newStructAlloca
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func InvokeConjOld(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
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
