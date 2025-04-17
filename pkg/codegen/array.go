package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareArrayInt(ir *ir.Module, arrayType *mTypes.ArrayTypeProps) {
	// 型: struct { i32* %arrElm, i64 %len}
	arrayIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	arrayIntType.SetName("array.int")
	ir.NewTypeDef("array.int", arrayIntType)
	arrayType.TypeInt = arrayIntType

}

func DeclareType(ir *ir.Module, arrayType *mTypes.ArrayTypeProps) {

	declareArrayInt(ir, arrayType)

	log.DebugMessage("array types declared")
}
