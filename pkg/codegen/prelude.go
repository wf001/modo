package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

var PreludeFunction = map[string]func(*Context, *mTypes.Node) value.Value{
	mTypes.PRELUDE_PRN:  PreludePrn,
	mTypes.PRELUDE_GET:  PreludeGet,
	mTypes.PRELUDE_CONJ: PreludeConj,
	//mTypes.LIB_CORE_ASSOC: PreludeAssoc,
	//mTypes.LIB_CORE_POP:   PreludePop,
	// nary
	mTypes.OPERATOR_ADD: PreludeAdd,
	mTypes.OPERATOR_SUB: PreludeSub,
	mTypes.OPERATOR_MUL: PreludeMul,
	mTypes.OPERATOR_DIV: PreludeDiv,
	mTypes.OPERATOR_EQ:  PreludeEq,
	mTypes.OPERATOR_GT:  PreludeGt,
	mTypes.OPERATOR_LT:  PreludeLt,
	mTypes.OPERATOR_AND: PreludeAnd,
	mTypes.OPERATOR_OR:  PreludeOr,
	// binary
	mTypes.OPERATOR_MOD: PreludeMod,
}

func declarePrelude(ir *ir.Module, prelude *mTypes.PreludeProps) {

	log.DebugMessage("vector types declared")
}
