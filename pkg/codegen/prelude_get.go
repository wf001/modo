package codegen

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeGet(ctx *Context, n *mTypes.Node) value.Value {
	i, _ := strconv.ParseInt(n.Next.Val, 10, 32)
	idx := constant.NewInt(types.I64, i)

	oldStructedVecPtr := n.IRValue
	structPtrType := oldStructedVecPtr.Type().(*types.PointerType)
	structType := structPtrType.ElemType.(*types.StructType)

	structedVecType, elemType := GetLLVMTypeFromString(structType.TypeName, ctx.prog.Prelude)
	nullPtr := constant.NewNull(types.NewPointer(elemType))

	oldStructedVec := ctx.block.NewLoad(structedVecType, oldStructedVecPtr)
	oldVecPtr := ctx.block.NewExtractValue(oldStructedVec, 0)
	oldLen := ctx.block.NewExtractValue(oldStructedVec, 1)
	maxIdx := ctx.block.NewSub(oldLen, constant.NewInt(types.I64, 1))

	exceedMaxIdx := ctx.block.NewICmp(enum.IPredSGT, idx, maxIdx)

	parent := ctx.block
	fn := ctx.function

	okBlock := fn.NewBlock("index.ok")
	failBlock := fn.NewBlock("index.fail")
	mergeBlock := fn.NewBlock("index.merge")

	parent.NewCondBr(exceedMaxIdx, failBlock, okBlock)

	// okBlock: 正常アクセス
	ctx.block = okBlock
	oldArrElemPtr := ctx.block.NewGetElementPtr(elemType, oldVecPtr, idx)
	oldElem := ctx.block.NewLoad(elemType, oldArrElemPtr)
	okPtr := ctx.block.NewAlloca(elemType)
	ctx.block.NewStore(oldElem, okPtr)
	ctx.block.NewBr(mergeBlock)

	// failBlock: 範囲外 → null
	ctx.block = failBlock
	ctx.block.NewBr(mergeBlock)

	// mergeBlock: φノードでブロックごとの戻り値を選択
	ctx.block = mergeBlock
	incs := []*ir.Incoming{
		ir.NewIncoming(okPtr, okBlock),
		ir.NewIncoming(nullPtr, failBlock),
	}
	result := ctx.block.NewPhi(incs...)

	return result
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func PreludeGetOld(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	value, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", value)
	}
	t := value.ElemType.(*types.ArrayType)

	i, _ := strconv.ParseInt(node.Next.Val, 10, 32)
	if i >= int64(t.Len) {
		log.Panic("Array index out of range: have %d but array length %d", i, node.Len)
	}
	elemPtr := block.NewGetElementPtr(
		value.ElemType,
		value,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, i),
	)
	elem := block.NewLoad(t.ElemType, elemPtr)

	return elem
}
