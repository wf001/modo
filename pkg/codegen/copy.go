package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

func CopyArray(
	block *ir.Block,
	libs *mTypes.Internal,
	allocSize *ir.InstMul, // The number of bytes to allocate
	oldArrType *types.ArrayType,
	oldArrPtr value.Value,
	newArrType *types.ArrayType,
) *ir.InstBitCast {
	allocatedPtr := block.NewCall(libs.Cstd.Malloc, allocSize)

	newArrPtr := block.NewBitCast(allocatedPtr, types.NewPointer(newArrType))

	for i := uint64(0); i < oldArrType.Len; i++ {
		elem := mTypes.LoadArrElem(block, oldArrPtr, oldArrType, i)
		destPtr := block.NewGetElementPtr(
			newArrType,
			newArrPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		block.NewStore(elem, destPtr)
	}
	return newArrPtr
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func CopyArrayOld(
	block *ir.Block,
	oldArrPtr *ir.InstAlloca, // Pointer to the allocated memory of the original array
	arrType types.Type, // The type of elements in the array (e.g., types.I32)
	newArrSize uint64, // The total size of the new array (number of elements to allocate)
	elemSize uint64, // The number of elements to copy from the old array to the new one
) *ir.InstAlloca {
	newArrType := types.NewArray(newArrSize, arrType)

	var newArr = block.NewAlloca(newArrType)
	for i := uint64(0); i < elemSize; i++ {
		oldElemPtr := block.NewGetElementPtr(
			oldArrPtr.ElemType,
			oldArrPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		// Arrayの要素の型取得
		elem := block.NewLoad(arrType, oldElemPtr)
		newElemPtr := block.NewGetElementPtr(
			newArrType,
			newArr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		block.NewStore(elem, newElemPtr)
	}
	return newArr

}
