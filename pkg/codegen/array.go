package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareVectorType(ir *ir.Module, vectorType *mTypes.VectorTypeProps) {
	// 型: struct { i32* %arrElm, i64 %len}
	vectorIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	vectorIntType.SetName("prelude.vector.int")
	ir.NewTypeDef("prelude.vector.int", vectorIntType)
	vectorType.TypeInt = vectorIntType

	// 型: struct { i8** %arrElm, i64 %len}
	arrayStringType := types.NewStruct(types.NewPointer(types.I8Ptr), types.I64)
	arrayStringType.SetName("prelude.vector.string")
	ir.NewTypeDef("prelude.vector.string", arrayStringType)
	vectorType.TypeString = arrayStringType

}

func DeclareType(ir *ir.Module, arrayType *mTypes.VectorTypeProps) {

	declareVectorType(ir, arrayType)

	log.DebugMessage("vector types declared")
}
