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
	{mTypes.OPERATOR_SUB, types.I32, types.I32}:     PreludeSub,
	{mTypes.OPERATOR_SUB, types.Float, types.Float}: PreludeFSub,
	{mTypes.OPERATOR_MUL, types.I32, types.I32}:     PreludeMul,
	{mTypes.OPERATOR_MUL, types.Float, types.Float}: PreludeFMul,
	{mTypes.OPERATOR_DIV, types.I32, types.I32}:     PreludeDiv,
	{mTypes.OPERATOR_DIV, types.Float, types.Float}: PreludeFDiv,
	{mTypes.OPERATOR_EQ, types.I1, types.I1}:        PreludeEq,
	{mTypes.OPERATOR_EQ, types.I32, types.I32}:      PreludeEq,
	{mTypes.OPERATOR_EQ, types.I8Ptr, types.I8Ptr}:  PreludeEq,
	{mTypes.OPERATOR_EQ, types.Float, types.Float}:  PreludeFEq,
	{mTypes.OPERATOR_GT, types.I32, types.I32}:      PreludeGt,
	{mTypes.OPERATOR_GT, types.Float, types.Float}:  PreludeFGt,
	{mTypes.OPERATOR_LT, types.I32, types.I32}:      PreludeLt,
	{mTypes.OPERATOR_LT, types.Float, types.Float}:  PreludeFLt,
	{mTypes.OPERATOR_MOD, types.I32, types.I32}:     PreludeMod,
	{mTypes.OPERATOR_MOD, types.Float, types.Float}: PreludeFMod,
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
	mTypes.PRELUDE_REDUCE: PreludeReduce,
	//mTypes.LIB_CORE_ASSOC: PreludeAssoc,
	//mTypes.LIB_CORE_POP:   PreludePop,
	mTypes.OPERATOR_ADD: OperatorDispatcher(mTypes.OPERATOR_ADD),
	mTypes.OPERATOR_SUB: OperatorDispatcher(mTypes.OPERATOR_SUB),
	mTypes.OPERATOR_MUL: OperatorDispatcher(mTypes.OPERATOR_MUL),
	mTypes.OPERATOR_DIV: OperatorDispatcher(mTypes.OPERATOR_DIV),
	mTypes.OPERATOR_EQ:  OperatorDispatcher(mTypes.OPERATOR_EQ),
	mTypes.OPERATOR_GT:  OperatorDispatcher(mTypes.OPERATOR_GT),
	mTypes.OPERATOR_LT:  OperatorDispatcher(mTypes.OPERATOR_LT),
	mTypes.OPERATOR_AND: PreludeAnd,
	mTypes.OPERATOR_OR:  PreludeOr,
	mTypes.OPERATOR_MOD: OperatorDispatcher(mTypes.OPERATOR_MOD),
}
