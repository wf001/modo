package lib

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
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true

	libs.Printf = &mTypes.BuiltinProp{
		FuncPtr: printfFunc,
	}

}

func declareCore(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
}
