package codegen

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeGet(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	var oldArrType *types.ArrayType

	var oldArrPtr value.Value

	switch v := node.IRValue.(type) {
	case *ir.InstCall, *ir.InstBitCast:
		oldArrType, _ = mTypes.AssertArrType(v)
		oldArrPtr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", node.IRValue)
	}

	i, _ := strconv.ParseInt(node.Next.Val, 10, 32)
	if i >= int64(oldArrType.Len) {
		log.Panic("Array index out of range: have %d but array length %d", i, node.Len)
	}

	newValue := block.NewAlloca(oldArrType.ElemType)

	loadedValue := mTypes.LoadArrElem(block, oldArrPtr, oldArrType, uint64(i))
	block.NewStore(loadedValue, newValue)

	return newValue
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
