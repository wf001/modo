package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	mTypes "github.com/wf001/modo/pkg/types"
)

func declarePrintf(
	module *ir.Module,
	internal *mTypes.Internal,
) {
	printfFunc := module.NewFunc(
		"printf",
		types.Void,
		ir.NewParam("format", types.I8Ptr),
	)
	printfFunc.Sig.Variadic = true

	internal.Cstd.Printf = printfFunc

}
func declareMalloc(
	module *ir.Module,
	libs *mTypes.Internal,
) {
	mallocFunc := module.NewFunc(
		"malloc",
		types.I8Ptr,
	)

	libs.Cstd.Malloc = mallocFunc

}
func declareMemcpy(
	module *ir.Module,
	libs *mTypes.Internal,
) {

	memcpyFunc := module.NewFunc(
		"llvm.memcpy.p0i8.p0i8.i64",
		types.Void,
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I64),
		ir.NewParam("", types.I1),
	)
	libs.Cstd.Memcpy = memcpyFunc
}
func declareStrcmp(
	module *ir.Module,
	libs *mTypes.Internal,
) {

	memcpyFunc := module.NewFunc(
		"strcmp",
		types.I1,
		ir.NewParam("", types.I8Ptr),
		ir.NewParam("", types.I8Ptr),
	)
	libs.Cstd.Strcmp = memcpyFunc
}

func declareCstd(ir *ir.Module, libs *mTypes.Internal) {
	declarePrintf(ir, libs)
	declareMalloc(ir, libs)
	declareMemcpy(ir, libs)
	declareStrcmp(ir, libs)
}
