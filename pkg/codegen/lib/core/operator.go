package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeAdd(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	})
}

func InvokeSub(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewSub(x, y)
	})
}

func InvokeMul(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewMul(x, y)
	})
}

func InvokeEq(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewICmp(enum.IPredEQ, x, y)
	})
}

func InvokeMod(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return block.NewSRem(node.IRValue, node.Next.IRValue)
}
