package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareArrayType(ir *ir.Module, arrayType *mTypes.ArrayTypeProps) {
	// 型: struct { i32* %arrElm, i64 %len}
	arrayIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	arrayIntType.SetName("array.int")
	ir.NewTypeDef("array.int", arrayIntType)
	arrayType.TypeInt = arrayIntType

	// 型: struct { i8** %arrElm, i64 %len}
	arrayStringType := types.NewStruct(types.NewPointer(types.I8Ptr), types.I64)
	arrayStringType.SetName("array.string")
	ir.NewTypeDef("array.string", arrayStringType)
	arrayType.TypeString = arrayStringType

}

func DeclareType(ir *ir.Module, arrayType *mTypes.ArrayTypeProps) {

	declareArrayType(ir, arrayType)

	log.DebugMessage("array types declared")
}
