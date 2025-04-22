package codegen

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func getVec(ctx *Context, n *mTypes.Node) value.Value {
	i, _ := strconv.ParseInt(n.Next.Val, 10, 32)
	idx := constant.NewInt(types.I64, i)

	structedVecPtr := n.IRValue
	structedVecType := mTypes.GetStructTypeFromPtr(structedVecPtr)
	elemType := structedVecType.Fields[0].(*types.PointerType).ElemType

	nullPtr := constant.NewNull(types.NewPointer(elemType))

	if i < 0 {
		return nullPtr
	}

	loadedStructedVec := ctx.block.NewLoad(structedVecType, structedVecPtr)
	loadedVecPtr := ctx.block.NewExtractValue(loadedStructedVec, 0)
	loadedLength := ctx.block.NewExtractValue(loadedStructedVec, 1)
	maxIdx := ctx.block.NewSub(loadedLength, constant.NewInt(types.I64, 1))

	isIdxOutOfRange := ctx.block.NewICmp(enum.IPredSGT, idx, maxIdx)
	isIdxOutOfRange.SetName(n.GetVarName("is.idx.out.of.range", ctx.block.Insts))

	inRangeBlock := ctx.function.NewBlock(n.GetBlockName("idx.in.range", ctx.function.Blocks))
	outOfRangeBlock := ctx.function.NewBlock(
		n.GetBlockName("idx.out.of.range", ctx.function.Blocks),
	)
	mergeBlock := ctx.function.NewBlock(n.GetBlockName("idx.merge", ctx.function.Blocks))

	ctx.block.NewCondBr(isIdxOutOfRange, outOfRangeBlock, inRangeBlock)

	// branched when specified index is in range
	ctx.block = inRangeBlock
	vecElemPtr := ctx.block.NewGetElementPtr(elemType, loadedVecPtr, idx)
	elem := ctx.block.NewLoad(elemType, vecElemPtr)
	okPtr := ctx.block.NewAlloca(elemType)
	ctx.block.NewStore(elem, okPtr)
	ctx.block.NewBr(mergeBlock)

	// branched when specified index is out of range
	ctx.block = outOfRangeBlock
	ctx.block.NewBr(mergeBlock)

	ctx.block = mergeBlock
	incs := []*ir.Incoming{
		ir.NewIncoming(okPtr, inRangeBlock),
		ir.NewIncoming(nullPtr, outOfRangeBlock),
	}
	result := ctx.block.NewPhi(incs...)

	return result
}

func getStruct(ctx *Context, n *mTypes.Node) value.Value {
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

	oldStruct := ctx.block.NewLoad(structeType, targetStructPtr)
	elemPtr := ctx.block.NewExtractValue(oldStruct, i)

	return elemPtr
}

func PreludeGet(ctx *Context, n *mTypes.Node) value.Value {
	tyPtr, okPtr := n.IRValue.Type().(*types.PointerType)
	tyStr, okStr := tyPtr.ElemType.(*types.StructType)

	if okPtr &&
		okStr &&
		ctx.prog.Declare.Type != nil &&
		ctx.prog.Declare.Type.Struct != nil {
		for k := range ctx.prog.Declare.Type.Struct {
			if k == tyStr.TypeName {
				return getStruct(ctx, n)
			}
		}
	}
	return getVec(ctx, n)
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func PreludeGetOld(block *ir.Block, internal *mTypes.Internal, node *mTypes.Node) value.Value {
	value, ok := node.IRValue.(*ir.InstAlloca)
	if !ok {
		log.Panic("Array elements must be ir.InstAlloca: have %+v", value)
	}
	t := value.ElemType.(*types.ArrayType)

	i, _ := strconv.ParseInt(node.Next.Val, 10, 32)
	if i >= int64(t.Len) {
		log.Panic("Array index out of range: have %d but array length %d", i, node.Len)
	}
	elemPtr := block.NewGetElementPtr(
		value.ElemType,
		value,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, i),
	)
	elem := block.NewLoad(t.ElemType, elemPtr)

	return elem
}
