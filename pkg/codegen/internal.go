package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareInternal(ir *ir.Module, libs *mTypes.Internal) {

	declareInternalConst(ir, libs)
	declareCstd(ir, libs)

	log.DebugMessage("built-in library declared")
}
