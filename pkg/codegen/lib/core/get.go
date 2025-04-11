package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeGet(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var formatStr *ir.Global

	for n := node; n != nil; n = n.Next {
		ty := n.IRValue.Type()

		if ty.Equal(types.I32) ||
			ty.Equal(types.I1) ||
			ty.Equal(types.I8Ptr) ||
			ty.Equal(types.Void) {
			PrnScalar(formatStr, libs, block, n)

		} else if t, ok := ty.(*types.VectorType); ok {
			PrnVector(libs, block, n, t)

		} else {
			log.Panic("unresolved type: have %+v", n)
		}

		if n.Next == nil {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatCR)
		} else {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}

	}
	return nil
}
