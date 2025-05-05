package codegen

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

type OpFunc func(*Context, *mTypes.Node) value.Value

type OpKey struct {
	Operator string
	Left     types.Type
	Right    types.Type
}

var OperatorDispatchTable = map[OpKey]OpFunc{
	{mTypes.OPERATOR_ADD, types.I32, types.I32}:     PreludeAdd,
	{mTypes.OPERATOR_ADD, types.Float, types.Float}: PreludeFAdd,
}

func OperatorDispatcher(op string) func(*Context, *mTypes.Node) value.Value {
	return func(ctx *Context, node *mTypes.Node) value.Value {
		if node.GetNodeSize() != 2 {
			log.Panic("%s: expects 2 arguments: %s", error.ERROR_SYNTAX_ERROR, op)
		}

		key := OpKey{op, node.IRValue.Type(), node.Next.IRValue.Type()}
		if fn, ok := OperatorDispatchTable[key]; ok {
			return fn(ctx, node)
		}
		log.Panic(
			"%s: arguments passed to operator must be same type",
			error.ERROR_UNDEFINE,
			op,
		)
		return nil
	}
}

var PreludeFunction = map[string]OpFunc{
	mTypes.PRELUDE_PRN:    PreludePrn,
	mTypes.PRELUDE_GET:    PreludeGet,
	mTypes.PRELUDE_NTH:    PreludeNth,
	mTypes.PRELUDE_CONJ:   PreludeConj,
	mTypes.PRELUDE_MAP:    PreludeMap,
	mTypes.PRELUDE_FILTER: PreludeFilter,
	//mTypes.LIB_CORE_ASSOC: PreludeAssoc,
	//mTypes.LIB_CORE_POP:   PreludePop,
	// nary
	mTypes.OPERATOR_ADD: OperatorDispatcher(mTypes.OPERATOR_ADD),
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
