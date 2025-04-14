package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

// arithmetic
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
func InvokeDiv(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewSDiv(x, y)
	})
}

func InvokeMod(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return block.NewSRem(node.IRValue, node.Next.IRValue)
}

// equality
func InvokeEq(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	// In case of comparing InstCall, compare the value itself by strcmp, not the address.
	_, ok := node.IRValue.(*ir.InstCall)
	if ok && node.Type == mTypes.TY_STR {
		fst := node.IRValue

		node = node.Next
		snd := node.IRValue

		cmpRes := block.NewCall(libs.Strcmp.FuncPtr, fst, snd)
		isEq := block.NewICmp(enum.IPredEQ, cmpRes, constant.NewInt(types.I1, 0))
		res := block.NewAnd(isEq, constant.NewInt(types.I1, 1))

		for n := node.Next; n != nil; n = n.Next {
			snd = n.IRValue
			cmpRes = block.NewCall(libs.Strcmp.FuncPtr, fst, snd)
			isEq = block.NewICmp(enum.IPredEQ, cmpRes, constant.NewInt(types.I1, 0))
			res = block.NewAnd(res, isEq)
		}
		return res

	} else {
		return invokeFold(node, func(x, y value.Value) value.Value {
			return block.NewICmp(enum.IPredEQ, x, y)
		})

	}
}
func InvokeGt(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewICmp(enum.IPredSGT, x, y)
	})
}
func InvokeLt(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewICmp(enum.IPredSLT, x, y)
	})
}

// logical
func InvokeAnd(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewAnd(x, y)
	})
}
func InvokeOr(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return invokeFold(node, func(x, y value.Value) value.Value {
		return block.NewOr(x, y)
	})
}
