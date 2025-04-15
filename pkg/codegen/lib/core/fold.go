package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	mTypes "github.com/wf001/modo/pkg/types"
)

func invokeFoldArithmetic(
	node *mTypes.Node,
	operate func(value.Value, value.Value) value.Value,
) value.Value {

	fst := node.IRValue

	node = node.Next
	snd := node.IRValue

	res := operate(fst, snd)

	for node = node.Next; node != nil; node = node.Next {
		snd := node.IRValue
		res = operate(res, snd)
	}
	return res
}
func invokeFoldPred(
	block *ir.Block,
	node *mTypes.Node,
	operate func(value.Value, value.Value) value.Value,
) value.Value {

	fst := node.IRValue

	node = node.Next
	snd := node.IRValue

	res := operate(fst, snd)

	for node = node.Next; node != nil; node = node.Next {
		snd := node.IRValue
		cmpRes := operate(fst, snd)
		res = block.NewAnd(cmpRes, res)
	}
	return res
}
