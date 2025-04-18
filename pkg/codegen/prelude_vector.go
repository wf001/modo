package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	mTypes "github.com/wf001/modo/pkg/types"
)

const PRELUDE_TYPENAME_VECINT = "prelude.vector.int"
const PRELUDE_TYPENAME_VECSTR = "prelude.vector.string"
const PRELUDE_TYPENAME_VECBOOL = "prelude.vector.bool"

func declareVectorType(ir *ir.Module, Prelude *mTypes.PreludeProps) {
	// i32 type vector
	// 型: struct { i32* %arrElm, i64 %len}
	vectorIntType := types.NewStruct(types.NewPointer(types.I32), types.I64)
	vectorIntType.SetName(PRELUDE_TYPENAME_VECINT)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECINT, vectorIntType)
	Prelude.Types.VectorInt = vectorIntType

	// i1 (bool) type vector
	// 型: struct { i1* %arrElm, i64 %len}
	vectorBoolType := types.NewStruct(types.NewPointer(types.I1), types.I64)
	vectorBoolType.SetName(PRELUDE_TYPENAME_VECBOOL)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECBOOL, vectorBoolType)
	Prelude.Types.VectorBool = vectorBoolType

	// i8* (string) type vector
	// 型: struct { i8** %arrElm, i64 %len}
	vectorStringType := types.NewStruct(types.NewPointer(types.I8Ptr), types.I64)
	vectorStringType.SetName(PRELUDE_TYPENAME_VECSTR)
	ir.NewTypeDef(PRELUDE_TYPENAME_VECSTR, vectorStringType)
	Prelude.Types.VectorString = vectorStringType

}

func GetLLVMTypeFromString(typeName string, prelude *mTypes.PreludeProps) (types.Type, types.Type) {

	switch typeName {

	case PRELUDE_TYPENAME_VECINT:
		return prelude.Types.VectorInt, types.I32

	case PRELUDE_TYPENAME_VECSTR:
		return prelude.Types.VectorString, types.I8Ptr

	case PRELUDE_TYPENAME_VECBOOL:
		return prelude.Types.VectorBool, types.I1
	}

	return nil, nil

}
