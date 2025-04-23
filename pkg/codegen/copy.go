package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	mTypes "github.com/wf001/modo/pkg/types"
)

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
			mTypes.I32zero,
			constant.NewInt(types.I32, int64(i)),
		)
		// Arrayの要素の型取得
		elem := block.NewLoad(arrType, oldElemPtr)
		newElemPtr := block.NewGetElementPtr(
			newArrType,
			newArr,
			mTypes.I32zero,
			constant.NewInt(types.I32, int64(i)),
		)
		block.NewStore(elem, newElemPtr)
	}
	return newArr

}
