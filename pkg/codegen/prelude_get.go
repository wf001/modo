package codegen

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeGet(ctx *Context, n *mTypes.Node) value.Value {
	tyPtr, _ := n.IRValue.Type().(*types.PointerType)
	tyStr, _ := tyPtr.ElemType.(*types.StructType)

	specifiedFields := ctx.prog.Declare.Type.Struct[tyStr.TypeName].Field[n.Next.Val]
	i := specifiedFields.Pos
	structeType, elemType := ctx.prog.Declare.Type.Struct[tyStr.TypeName].Types, specifiedFields.Type

	targetStructPtr := n.IRValue

	if structeType == nil || elemType == nil {
		// Note: is it ok returing i32* anytime?
		return constant.NewNull(types.NewPointer(types.I32))
	}

	loadedStruct := ctx.block.NewLoad(structeType, targetStructPtr)
	elemPtr := ctx.block.NewExtractValue(loadedStruct, i)

	return elemPtr
}
