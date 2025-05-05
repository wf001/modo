package codegen

import (
	"github.com/llir/llvm/ir"

	mTypes "github.com/wf001/modo/pkg/types"
)

func (ctx Context) NewBlock(name string, n *mTypes.Node) *ir.Block {
	return ctx.function.NewBlock(n.GetBlockName(name, ctx.function.Blocks))
}
