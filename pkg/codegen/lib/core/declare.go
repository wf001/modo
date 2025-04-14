package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	mTypes "github.com/wf001/modo/pkg/types"
)

func declarePrintf(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {
	printfFunc := module.NewFunc(
		"printf",
		types.Void,
		ir.NewParam("format", types.I8Ptr),
	)
	printfFunc.Sig.Variadic = true

	libs.Printf = &mTypes.BuiltinProp{
		FuncPtr: printfFunc,
	}

}
func declareMalloc(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {
	mallocFunc := module.NewFunc(
		"malloc",
		types.I8Ptr,
	)

	libs.Malloc = &mTypes.BuiltinProp{
		FuncPtr: mallocFunc,
	}

}
func declareMemcpy(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {

	memcpyFunc := module.NewFunc(
		"llvm.memcpy.p0i8.p0i8.i64",
		types.Void,
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I64),
		ir.NewParam("", types.I1),
	)
	libs.Memcpy = &mTypes.BuiltinProp{
		FuncPtr: memcpyFunc,
	}
}
func declareStrcmp(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {

	memcpyFunc := module.NewFunc(
		"strcmp",
		types.I1,
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I8Ptr),
	)
	libs.Strcmp = &mTypes.BuiltinProp{
		FuncPtr: memcpyFunc,
	}
}

func Declare(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
	declareMalloc(ir, libs)
	declareMemcpy(ir, libs)
	declareStrcmp(ir, libs)
}
