package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func PreludeAssoc(ctx *Context, n *mTypes.Node) value.Value {
	srcPtr := n.IRValue

	tyPtr, isTyPtr := srcPtr.Type().(*types.PointerType)
	tyStr, isTyStr := tyPtr.ElemType.(*types.StructType)
	if !isTyPtr || !isTyStr {
		return constant.NewNull(types.NewPointer(types.I32))
	}

	typeInfo := ctx.prog.Declare.Type.Struct[tyStr.TypeName]
	structType := typeInfo.Types

	nullStructPtr := constant.NewNull(types.NewPointer(structType))
	gepEndPtr := ctx.block.NewGetElementPtr(structType, nullStructPtr, mTypes.I32one)
	mallocSize := ctx.block.NewPtrToInt(gepEndPtr, types.I64)
	rawPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, mallocSize)
	destPtr := ctx.block.NewBitCast(rawPtr, types.NewPointer(structType))

	// deep copy
	for _, field := range typeInfo.Field {
		fieldIndex := field.Pos
		fieldType := field.Type

		origFieldPtr := ctx.block.NewGetElementPtr(
			structType,
			srcPtr,
			mTypes.I32zero,
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
		origVal := ctx.block.NewLoad(fieldType, origFieldPtr)

		copyFieldPtr := ctx.block.NewGetElementPtr(
			structType,
			destPtr,
			mTypes.I32zero,
			constant.NewInt(types.I32, int64(fieldIndex)),
		)

		if ptrType, ok := fieldType.(*types.PointerType); ok {
			isNull := ctx.block.NewICmp(enum.IPredEQ, origVal, constant.NewNull(ptrType))

			nullBlock := ctx.NewBlock("copy.null", n)
			nonNullBlock := ctx.NewBlock("copy.non.null", n)
			mergeBlock := ctx.NewBlock("copy.merge", n)

			ctx.block.NewCondBr(isNull, nullBlock, nonNullBlock)

			nullBlock.NewStore(constant.NewNull(ptrType), copyFieldPtr)
			nullBlock.NewBr(mergeBlock)

			elemType := ptrType.ElemType
			nullElemPtr := constant.NewNull(types.NewPointer(elemType))
			gepEnd := nonNullBlock.NewGetElementPtr(elemType, nullElemPtr, mTypes.I32one)
			allocSize := nonNullBlock.NewPtrToInt(gepEnd, types.I64)

			malloc := nonNullBlock.NewCall(ctx.internal.Cstd.Malloc, allocSize)
			newPtr := nonNullBlock.NewBitCast(malloc, ptrType)

			loadedElem := nonNullBlock.NewLoad(elemType, origVal)
			nonNullBlock.NewStore(loadedElem, newPtr)

			nonNullBlock.NewStore(newPtr, copyFieldPtr)
			nonNullBlock.NewBr(mergeBlock)

			ctx.block = mergeBlock

		} else {
			ctx.block.NewStore(origVal, copyFieldPtr)
		}
	}

	fieldName := n.Next.Val
	fieldInfo := typeInfo.Field[fieldName]
	fieldIndex := fieldInfo.Pos
	fieldType := fieldInfo.Type

	fieldPtr := ctx.block.NewGetElementPtr(
		structType,
		destPtr,
		mTypes.I32zero,
		constant.NewInt(types.I32, int64(fieldIndex)),
	)

	val := n.Next.Next.IRValue

	if _, ok := val.(*ir.InstBitCast); ok {
		if ptrType, ok := fieldType.(*types.PointerType); ok {
			nullVal := constant.NewNull(ptrType)
			ctx.block.NewStore(nullVal, fieldPtr)
			ctx.block.NewStore(val, fieldPtr)
			return destPtr
		}
	}

	ctx.block.NewStore(val, fieldPtr)
	return destPtr
}
