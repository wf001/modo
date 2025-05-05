package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

// arithmetic
func PreludeAdd(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewAdd(node.IRValue, node.Next.IRValue)
}

func PreludeFAdd(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFAdd(node.IRValue, node.Next.IRValue)
}

func PreludeSub(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewSub(node.IRValue, node.Next.IRValue)
}
func PreludeFSub(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFSub(node.IRValue, node.Next.IRValue)
}

func PreludeMul(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewMul(node.IRValue, node.Next.IRValue)
}

func PreludeFMul(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFMul(node.IRValue, node.Next.IRValue)
}

func PreludeDiv(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewSDiv(node.IRValue, node.Next.IRValue)
}

func PreludeFDiv(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFDiv(node.IRValue, node.Next.IRValue)
}

func PreludeMod(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewSRem(node.IRValue, node.Next.IRValue)
}
func PreludeFMod(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFRem(node.IRValue, node.Next.IRValue)
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
		isEq := ctx.block.NewICmp(enum.IPredEQ, cmpRes, mTypes.I1zero)
		return ctx.block.NewAnd(isEq, mTypes.I1one)

	} else {
		return ctx.block.NewICmp(enum.IPredEQ, node.IRValue, node.Next.IRValue)
	}
}

func PreludeFEq(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFCmp(enum.FPredOEQ, node.IRValue, node.Next.IRValue)
}

func PreludeGt(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewICmp(enum.IPredSGT, node.IRValue, node.Next.IRValue)
}

func PreludeFGt(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFCmp(enum.FPredOGT, node.IRValue, node.Next.IRValue)
}

func PreludeLt(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewICmp(enum.IPredSLT, node.IRValue, node.Next.IRValue)
}

func PreludeFLt(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewFCmp(enum.FPredOLT, node.IRValue, node.Next.IRValue)
}

// logical
func PreludeAnd(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewAnd(node.IRValue, node.Next.IRValue)
}
func PreludeOr(ctx *Context, node *mTypes.Node) value.Value {
	return ctx.block.NewOr(node.IRValue, node.Next.IRValue)
}
