package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func prnScalar(
	ctx *Context,
	n *mTypes.Node,
) {
	value := n.IRValue
	rootTy := value.Type()
	formatStr, _ := mTypes.GetPrintFormat(rootTy, ctx.internal)

	if rootTy.Equal(types.I1) {
		value = ctx.block.NewSelect(
			value,
			ctx.internal.GlobalConst.StringTrue,
			ctx.internal.GlobalConst.StringFalse,
		)
		ctx.block.NewCall(ctx.internal.Cstd.Printf, formatStr, value)

	} else if rootTy.Equal(types.Void) {
		value = ctx.internal.GlobalConst.StringNil
		ctx.block.NewCall(ctx.internal.Cstd.Printf, formatStr, value)

	} else if pointerElemTy, isPtr := rootTy.(*types.PointerType); isPtr {

		if isStr := rootTy.Equal(types.I8Ptr); isStr {
			ctx.block.NewCall(ctx.internal.Cstd.Printf, formatStr, value)

		} else {
			isNull := ctx.block.NewICmp(
				enum.IPredEQ,
				value,
				constant.NewNull(types.NewPointer(pointerElemTy)),
			)

			// if value is null ptr
			nullBlock := ctx.function.NewBlock(n.GetBlockName("print.null.ptr", ctx.function.Blocks))
			nullBlock.NewCall(ctx.internal.Cstd.Printf, ctx.internal.GlobalConst.FormatStr, ctx.internal.GlobalConst.StringNil)

			// if value is not null ptr
			nonNullBlock := ctx.function.NewBlock(n.GetBlockName("print.non.null", ctx.function.Blocks))
			formatStr, _ = mTypes.GetPrintFormat(pointerElemTy.ElemType, ctx.internal)

			ptr := nonNullBlock.NewLoad(pointerElemTy, value)
			value = nonNullBlock.NewGetElementPtr(
				pointerElemTy.ElemType,
				ptr,
				constant.NewInt(types.I32, 0),
			)

			nonNullBlock.NewCall(ctx.internal.Cstd.Printf, formatStr, value)

			endBlock := ctx.function.NewBlock(n.GetBlockName("print.end", ctx.function.Blocks))
			nullBlock.NewBr(endBlock)
			nonNullBlock.NewBr(endBlock)

			ctx.block.NewCondBr(isNull, nullBlock, nonNullBlock)
			ctx.block = endBlock
		}
	} else {
		// Except string, bool, nil, pointer type value
		ctx.block.NewCall(ctx.internal.Cstd.Printf, formatStr, value)
	}

}

// Note: NOT USED
// see newVectorGlobal comment
func prnVector(internal *mTypes.Internal, block *ir.Block, n *mTypes.Node) {
	var arr value.Value
	var arrType *types.ArrayType

	switch v := n.IRValue.(type) {
	// The vector on global space: InstCall,
	// The vector on local space: InstBitCast
	case *ir.InstCall, *ir.InstBitCast:
		arrType, _ = mTypes.AssertArrType(v)
		arr = v

	default:
		log.Panic("Unsupported IRValue type: %#+v", n.IRValue)
	}

	block.NewCall(internal.Cstd.Printf, internal.GlobalConst.StringBracketOpen)

	for i := uint64(0); i < arrType.Len; i++ {
		elem := mTypes.LoadArrElem(block, arr, arrType, i)

		formatStr, _ := mTypes.GetPrintFormat(elem.ElemType, internal)
		if arrType.ElemType == types.I1 {
			v := block.NewSelect(
				elem,
				internal.GlobalConst.StringTrue,
				internal.GlobalConst.StringFalse,
			)
			block.NewCall(internal.Cstd.Printf, formatStr, v)
		} else {
			block.NewCall(internal.Cstd.Printf, formatStr, elem)
		}

		if i < uint64(arrType.Len-1) {
			block.NewCall(internal.Cstd.Printf, internal.GlobalConst.StringComma)
			block.NewCall(internal.Cstd.Printf, internal.GlobalConst.StringSpace)
		}
	}

	block.NewCall(internal.Cstd.Printf, internal.GlobalConst.StringBracketClose)
}

func prnStructVector(
	ctx *Context,
	n *mTypes.Node,
) {
	v := n.IRValue

	ty := getStructTypeFromPtr(v)
	elemTy := ty.Fields[0].(*types.PointerType).ElemType

	isNull := ctx.block.NewICmp(
		enum.IPredEQ,
		v,
		constant.NewNull(types.NewPointer(ty)),
	)

	// if value is null ptr
	nullBlock := ctx.function.NewBlock(n.GetBlockName("print.null.ptr", ctx.function.Blocks))
	nullBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.FormatStr,
		ctx.internal.GlobalConst.StringNil,
	)

	// if value is not null ptr
	nonNullBlock := ctx.function.NewBlock(n.GetBlockName("print.non.null", ctx.function.Blocks))

	formatStr, _ := mTypes.GetPrintFormat(elemTy, ctx.internal)
	loaded := nonNullBlock.NewLoad(ty, n.IRValue)

	// 構造体のフィールドから arrPtr と len を取り出す
	resultArrPtr := nonNullBlock.NewExtractValue(loaded, 0)
	resultLen := nonNullBlock.NewExtractValue(loaded, 1)

	// インデックスの初期化
	idx := nonNullBlock.NewAlloca(types.I64)
	idx.SetName(n.GetVarName("idx", nonNullBlock.Insts))
	nonNullBlock.NewStore(constant.NewInt(types.I64, 0), idx)

	loopBlock := ctx.function.NewBlock(n.GetBlockName("printf.loop", ctx.function.Blocks))
	continueBlock := ctx.function.NewBlock(n.GetBlockName("printf.continue", ctx.function.Blocks))
	loopEndBlock := ctx.function.NewBlock(n.GetBlockName("printf.loop.end", ctx.function.Blocks))

	nonNullBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.StringBracketOpen,
	)
	nonNullBlock.NewBr(loopBlock)

	// ループ内部の処理
	// i をロード
	i := loopBlock.NewLoad(types.I64, idx)

	// 配列の各要素を取り出して表示
	elemPtr := loopBlock.NewGetElementPtr(elemTy, resultArrPtr, i)
	elem := loopBlock.NewLoad(elemTy, elemPtr)

	if elem.ElemType == types.I1 {
		boolValueStr := loopBlock.NewSelect(
			elem,
			ctx.internal.GlobalConst.StringTrue,
			ctx.internal.GlobalConst.StringFalse,
		)
		loopBlock.NewCall(ctx.internal.Cstd.Printf, formatStr, boolValueStr)
	} else {
		loopBlock.NewCall(
			ctx.internal.Cstd.Printf,
			formatStr,
			elem,
		)

	}

	// i++
	nextI := loopBlock.NewAdd(i, constant.NewInt(types.I64, 1))
	loopBlock.NewStore(nextI, idx)

	continueBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.StringComma,
	)
	continueBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.StringSpace,
	)
	continueBlock.NewBr(loopBlock)

	// i < len ?
	cond := loopBlock.NewICmp(enum.IPredSLT, nextI, resultLen)
	loopBlock.NewCondBr(cond, continueBlock, loopEndBlock)
	loopEndBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.StringBracketClose,
	)

	endBlock := ctx.function.NewBlock(n.GetBlockName("print.end", ctx.function.Blocks))
	nullBlock.NewBr(endBlock)
	loopEndBlock.NewBr(endBlock)

	ctx.block.NewCondBr(isNull, nullBlock, nonNullBlock)
	ctx.block = endBlock
}

func PreludePrn(
	ctx *Context,
	node *mTypes.Node,
) value.Value {

	for n := node; n != nil; n = n.Next {
		if mTypes.IsScalar(n.IRValue) {
			prnScalar(ctx, n)

		} else if p, ok := n.IRValue.Type().(*types.PointerType); ok {
			if _, ok := p.ElemType.(*types.StructType); ok {
				prnStructVector(ctx, n)
			} else {
				prnScalar(ctx, n)
			}
		} else if _, ok := n.IRValue.Type().(*types.StructType); ok {
			prnStructVector(ctx, n)

		} else {
			log.Panic("unresolved type: have %+v", n)
		}

		if n.Next == nil {
			ctx.block.NewCall(
				ctx.internal.Cstd.Printf,
				ctx.internal.GlobalConst.StringCR,
			)
		} else {
			ctx.block.NewCall(ctx.internal.Cstd.Printf, ctx.internal.GlobalConst.StringSpace)
		}
	}

	return nil
}
