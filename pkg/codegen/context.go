package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"

	mTypes "github.com/wf001/modo/pkg/types"
)

func (ctx Context) NewBlock(name string, n *mTypes.Node) *ir.Block {
	return ctx.function.NewBlock(ctx.GetBlockName(name, n))
}

func (ctx Context) GetBlockName(s string, node *mTypes.Node) string {
	return fmt.Sprintf("%s.%p.%d", s, node, len(ctx.block.Insts))
}
