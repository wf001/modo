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
	libs *mTypes.Internal,
	block *ir.Block,
	n *mTypes.Node,
) {
	value := n.IRValue
	rootTy := value.Type()
	formatStr, _ := mTypes.GetPrintFormat(rootTy, libs)

	if rootTy.Equal(types.I1) {
		value = block.NewSelect(
			value,
			libs.GlobalConst.StringTrue,
			libs.GlobalConst.StringFalse,
		)

	} else if rootTy.Equal(types.Void) {
		value = libs.GlobalConst.StringNil

	} else if pointerElemTy, isPtr := rootTy.(*types.PointerType); isPtr {
		if isStr := rootTy.Equal(types.I8Ptr); !isStr {
			ptr := block.NewLoad(pointerElemTy, value)
			value = block.NewGetElementPtr(
				pointerElemTy,
				ptr,
				constant.NewInt(types.I32, 0),
			)
			formatStr, _ = mTypes.GetPrintFormat(pointerElemTy.ElemType, libs)
		}
	}

	block.NewCall(libs.Cstd.Printf, formatStr, value)
}

// Note: NOT USED
// see newVectorGlobal comment
func prnVector(libs *mTypes.Internal, block *ir.Block, n *mTypes.Node) {
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

	block.NewCall(libs.Cstd.Printf, libs.GlobalConst.StringBracketOpen)

	for i := uint64(0); i < arrType.Len; i++ {
		elem := mTypes.LoadArrElem(block, arr, arrType, i)

		formatStr, _ := mTypes.GetPrintFormat(elem.ElemType, libs)
		if arrType.ElemType == types.I1 {
			v := block.NewSelect(
				elem,
				libs.GlobalConst.StringTrue,
				libs.GlobalConst.StringFalse,
			)
			block.NewCall(libs.Cstd.Printf, formatStr, v)
		} else {
			block.NewCall(libs.Cstd.Printf, formatStr, elem)
		}

		if i < uint64(arrType.Len-1) {
			block.NewCall(libs.Cstd.Printf, libs.GlobalConst.StringComma)
			block.NewCall(libs.Cstd.Printf, libs.GlobalConst.StringSpace)
		}
	}

	block.NewCall(libs.Cstd.Printf, libs.GlobalConst.StringBracketClose)
}

func prnStructVector(
	ctx *Context,
	n *mTypes.Node,
) {
	v := n.IRValue
	structPtrType, _ := v.Type().(*types.PointerType)
	structType := structPtrType.ElemType.(*types.StructType)
	var m1 = map[string]*types.StructType{
		"prelude.vector.int":    ctx.prog.Prelude.Types.VectorInt,
		"prelude.vector.string": ctx.prog.Prelude.Types.VectorString,
	}
	var m2 = map[string]types.Type{
		"prelude.vector.int":    types.I32,
		"prelude.vector.string": types.I8Ptr,
	}
	var m3 = map[string]*ir.Global{
		"prelude.vector.int":    ctx.prog.Internal.GlobalConst.FormatDigit,
		"prelude.vector.string": ctx.prog.Internal.GlobalConst.FormatStr,
	}
	ty := m1[structType.TypeName]
	elemTy := m2[structType.TypeName]
	formatStr := m3[structType.TypeName]
	loaded := ctx.block.NewLoad(ty, n.IRValue)

	// 構造体のフィールドから arrPtr と len を取り出す
	resultArrPtr := ctx.block.NewExtractValue(loaded, 0)
	resultLen := ctx.block.NewExtractValue(loaded, 1)

	// インデックスの初期化
	idx := ctx.block.NewAlloca(types.I64)
	idx.SetName(n.GetVarName("idx", ctx.block.Insts))
	ctx.block.NewStore(constant.NewInt(types.I64, 0), idx)

	loopBlock := ctx.function.NewBlock(n.GetBlockName("printf.loop", ctx.function.Blocks))
	continueBlock := ctx.function.NewBlock(n.GetBlockName("printf.continue", ctx.function.Blocks))
	endBlock := ctx.function.NewBlock(n.GetBlockName("printf.end", ctx.function.Blocks))

	ctx.block.NewCall(
		ctx.prog.Internal.Cstd.Printf,
		ctx.prog.Internal.GlobalConst.StringBracketOpen,
	)
	ctx.block.NewBr(loopBlock)

	// ループ内部の処理
	// i をロード
	i := loopBlock.NewLoad(types.I64, idx)

	// 配列の各要素を取り出して表示
	elemPtr := loopBlock.NewGetElementPtr(elemTy, resultArrPtr, i)
	elem := loopBlock.NewLoad(elemTy, elemPtr)
	loopBlock.NewCall(
		ctx.prog.Internal.Cstd.Printf,
		formatStr,
		elem,
	)

	// i++
	nextI := loopBlock.NewAdd(i, constant.NewInt(types.I64, 1))
	loopBlock.NewStore(nextI, idx)

	continueBlock.NewCall(
		ctx.prog.Internal.Cstd.Printf,
		ctx.prog.Internal.GlobalConst.StringComma,
	)
	continueBlock.NewCall(
		ctx.prog.Internal.Cstd.Printf,
		ctx.prog.Internal.GlobalConst.StringSpace,
	)
	continueBlock.NewBr(loopBlock)

	// i < len ?
	cond := loopBlock.NewICmp(enum.IPredSLT, nextI, resultLen)
	loopBlock.NewCondBr(cond, continueBlock, endBlock)
	ctx.block = endBlock
	ctx.block.NewCall(
		ctx.prog.Internal.Cstd.Printf,
		ctx.prog.Internal.GlobalConst.StringBracketClose,
	)

}

func PreludePrn(
	ctx *Context,
	node *mTypes.Node,
) value.Value {

	for n := node; n != nil; n = n.Next {
		if mTypes.IsScalar(n.IRValue) {
			prnScalar(ctx.prog.Internal, ctx.block, n)

		} else if _, ok := n.IRValue.Type().(*types.PointerType); ok {
			prnStructVector(ctx, n)

		} else if _, ok := n.IRValue.Type().(*types.StructType); ok {
			prnStructVector(ctx, n)

		} else {
			log.Panic("unresolved type: have %+v", n)
		}

		if n.Next == nil {
			ctx.block.NewCall(
				ctx.prog.Internal.Cstd.Printf,
				ctx.prog.Internal.GlobalConst.StringCR,
			)
		} else {
			ctx.block.NewCall(ctx.prog.Internal.Cstd.Printf, ctx.prog.Internal.GlobalConst.StringSpace)
		}
	}

	return nil
}
