package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeAssoc(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	value, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", value)
	}
	t := value.ElemType.(*types.ArrayType)
	arrType := types.NewArray(t.Len, t.ElemType)
	var newArr = block.NewAlloca(arrType)

	for i := uint64(0); i < t.Len; i++ {
		oldElemPtr := block.NewGetElementPtr(
			value.ElemType,
			value,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		// Arrayの要素の型取得
		elem := block.NewLoad(t.ElemType, oldElemPtr)
		newElemPtr := block.NewGetElementPtr(
			arrType,
			newArr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		block.NewStore(elem, newElemPtr)
	}

	newElemPtr := block.NewGetElementPtr(
		arrType,
		newArr,
		constant.NewInt(types.I32, 0),
		node.Next.IRValue,
	)
	block.NewStore(node.Next.Next.IRValue, newElemPtr)

	return newArr
}
