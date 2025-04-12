package vector

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func CopyArray(
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
