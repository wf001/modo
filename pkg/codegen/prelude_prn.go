package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func genNilBlock(
	ctx *Context,
	elemType *types.PointerType,
	n *mTypes.Node,
	v value.Value,
) (*ir.Block, *ir.Block, *ir.Block) {
	isNull := ctx.block.NewICmp(
		enum.IPredEQ,
		v,
		constant.NewNull(elemType),
	)

	// if value is null ptr
	nullBlock := ctx.NewBlock("prn.null", n)
	nullBlock.NewCall(
		ctx.internal.Cstd.Printf,
		ctx.internal.GlobalConst.FormatStr,
		ctx.internal.GlobalConst.StringNil,
	)

	// if value is not null ptr
	nonNullBlock := ctx.NewBlock("prn.non.null", n)

	endBlock := ctx.NewBlock("prn.exit", n)
	nullBlock.NewBr(endBlock)
	nonNullBlock.NewBr(endBlock)

	ctx.block.NewCondBr(isNull, nullBlock, nonNullBlock)

	return nullBlock, nonNullBlock, endBlock
}

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
			_, nonNullBlock, endBlock := genNilBlock(ctx, types.NewPointer(pointerElemTy), n, value)

			formatStr, _ = mTypes.GetPrintFormat(pointerElemTy.ElemType, ctx.internal)

			ptr := nonNullBlock.NewLoad(pointerElemTy, value)
			value = nonNullBlock.NewGetElementPtr(
				pointerElemTy.ElemType,
				ptr,
				mTypes.I32zero,
			)

			nonNullBlock.NewCall(ctx.internal.Cstd.Printf, formatStr, value)

			ctx.block = endBlock
		}
	} else {
		// Except string, bool, nil, pointer type value
		ctx.block.NewCall(ctx.internal.Cstd.Printf, formatStr, value)
	}

}

func prnStructVector(
	ctx *Context,
	n *mTypes.Node,
) {
	v := n.IRValue

	ty := mTypes.GetStructTypeFromPtr(v)
	elemTy := ty.Fields[0].(*types.PointerType).ElemType

	_, nonNullBlock, endBlock := genNilBlock(ctx, types.NewPointer(ty), n, v)

	// if value is not null ptr

	formatStr, _ := mTypes.GetPrintFormat(elemTy, ctx.internal)
	loaded := nonNullBlock.NewLoad(ty, n.IRValue)

	// 構造体のフィールドから arrPtr と len を取り出す
	resultArrPtr := nonNullBlock.NewExtractValue(loaded, 0)
	resultLen := nonNullBlock.NewExtractValue(loaded, 1)

	// インデックスの初期化
	idx := nonNullBlock.NewAlloca(types.I64)
	idx.SetName(n.GetVarName("prn.cur.idx", nonNullBlock.Insts))
	nonNullBlock.NewStore(mTypes.I64zero, idx)

	loopBlock := ctx.NewBlock("prn.vec.loop.enter", n)
	continueBlock := ctx.NewBlock("prn.vec.continue", n)
	loopEndBlock := ctx.NewBlock("prn.vec.loop.exit", n)

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
	nextI := loopBlock.NewAdd(i, mTypes.I64one)
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

	loopEndBlock.NewBr(endBlock)

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
			log.Panic("%s: unexpected value passed to prn function: have %+v", error.ERROR_UNDEFINE, n)
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
