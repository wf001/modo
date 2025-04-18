package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	mTypes "github.com/wf001/modo/pkg/types"
)

func declareVectorType(ir *ir.Module, Prelude *mTypes.PreludeProps) {
	// 型: struct { i32* %arrElm, i64 %len}
	vectorIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	vectorIntType.SetName("prelude.vector.int")
	ir.NewTypeDef("prelude.vector.int", vectorIntType)
	Prelude.Types.VectorInt = vectorIntType

	// 型: struct { i8** %arrElm, i64 %len}
	vectorStringType := types.NewStruct(types.NewPointer(types.I8Ptr), types.I64)
	vectorStringType.SetName("prelude.vector.string")
	ir.NewTypeDef("prelude.vector.string", vectorStringType)
	Prelude.Types.VectorString = vectorStringType

}

func GetLLVMTypeFromString(typeName string, prelude *mTypes.PreludeProps) (types.Type, types.Type) {

	switch typeName {

	case "prelude.vector.int":
		return prelude.Types.VectorInt, types.I32
	case "prelude.vector.string":
		return prelude.Types.VectorString, types.I8Ptr
	}
	return nil, nil

}
