package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareInternal(ir *ir.Module, internal *mTypes.Internal) {

	declareInternalConst(ir, internal)
	declareCstd(ir, internal)

	log.DebugMessage("declared internal global constant variable")
}
