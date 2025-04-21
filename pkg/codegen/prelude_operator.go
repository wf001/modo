package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

// arithmetic
func PreludeAdd(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldArithmetic(
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewAdd(x, y)
		})
}

func PreludeSub(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldArithmetic(
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewSub(x, y)
		})
}

func PreludeMul(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldArithmetic(
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewMul(x, y)
		})
}
func PreludeDiv(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldArithmetic(
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewSDiv(x, y)
		})
}

func PreludeMod(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewSRem(node.IRValue, node.Next.IRValue)
}

// equality
func PreludeEq(ctx *Context, node *mTypes.Node) value.Value {
	// In case of comparing InstCall, compare the value itself by strcmp, not the address.
	_, ok := node.IRValue.(*ir.InstCall)
	if ok && node.Type != nil && node.Type.Value == mTypes.TY_STR {
		fst := node.IRValue

		node = node.Next
		snd := node.IRValue

		cmpRes := ctx.block.NewCall(ctx.internal.Cstd.Strcmp, fst, snd)
		isEq := ctx.block.NewICmp(enum.IPredEQ, cmpRes, constant.NewInt(types.I1, 0))
		res := ctx.block.NewAnd(isEq, constant.NewInt(types.I1, 1))

		for n := node.Next; n != nil; n = n.Next {
			snd = n.IRValue
			cmpRes = ctx.block.NewCall(ctx.internal.Cstd.Strcmp, fst, snd)
			isEq = ctx.block.NewICmp(enum.IPredEQ, cmpRes, constant.NewInt(types.I1, 0))
			res = ctx.block.NewAnd(res, isEq)
		}
		return res

	} else {
		return invokeFoldPred(
			ctx.block,
			node,
			func(x, y value.Value) value.Value {
				return ctx.block.NewICmp(enum.IPredEQ, x, y)
			})

	}
}
func PreludeGt(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldPred(
		ctx.block,
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewICmp(enum.IPredSGT, x, y)
		})
}
func PreludeLt(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldPred(
		ctx.block,
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewICmp(enum.IPredSLT, x, y)
		})
}

// logical
func PreludeAnd(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldPred(
		ctx.block,
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewAnd(x, y)
		})
}
func PreludeOr(ctx *Context, node *mTypes.Node) value.Value {
	return invokeFoldPred(
		ctx.block,
		node,
		func(x, y value.Value) value.Value {
			return ctx.block.NewOr(x, y)
		})
}
