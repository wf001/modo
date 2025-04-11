package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, *mTypes.Node) value.Value{
	mTypes.LIB_CORE_PRN:   InvokePrn,
	mTypes.LIB_CORE_GET:   InvokeGet,
	mTypes.LIB_CORE_CONJ:  InvokeConj,
	mTypes.LIB_CORE_ASSOC: InvokeAssoc,
	mTypes.LIB_CORE_POP:   InvokePop,
	// nary
	mTypes.OPERATOR_ADD: InvokeAdd,
	mTypes.OPERATOR_SUB: InvokeSub,
	mTypes.OPERATOR_MUL: InvokeMul,
	mTypes.OPERATOR_DIV: InvokeDiv,
	mTypes.OPERATOR_EQ:  InvokeEq,
	mTypes.OPERATOR_GT:  InvokeGt,
	mTypes.OPERATOR_LT:  InvokeLt,
	mTypes.OPERATOR_AND: InvokeAnd,
	mTypes.OPERATOR_OR:  InvokeOr,
	// binary
	mTypes.OPERATOR_MOD: InvokeMod,
}
