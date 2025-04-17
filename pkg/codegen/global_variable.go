package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareVariable(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	globalVars := &mTypes.BuiltinGlobalVarsProp{}

	globalVars.FormatDigit = ir.NewGlobalDef(
		"format.digit",
		constant.NewCharArrayFromString("%d\x00"),
	)
	globalVars.FormatStr = ir.NewGlobalDef(
		"format.string",
		constant.NewCharArrayFromString("%s\x00"),
	)
	globalVars.FormatSpace = ir.NewGlobalDef(
		"format.space",
		constant.NewCharArrayFromString(" \x00"),
	)
	globalVars.FormatCR = ir.NewGlobalDef(
		"format.cr",
		constant.NewCharArrayFromString("\n\x00"),
	)
	globalVars.TrueValue = ir.NewGlobalDef(
		"value.true",
		constant.NewCharArrayFromString("true\x00"),
	)
	globalVars.FalseValue = ir.NewGlobalDef(
		"value.false",
		constant.NewCharArrayFromString("false\x00"),
	)
	globalVars.NilValue = ir.NewGlobalDef(
		"value.nil",
		constant.NewCharArrayFromString("nil\x00"),
	)
	globalVars.FormatBracketOpen = ir.NewGlobalDef(
		"format.bracket.open",
		constant.NewCharArrayFromString("[\x00"),
	)
	globalVars.FormatBracketClose = ir.NewGlobalDef(
		"format.bracket.close",
		constant.NewCharArrayFromString("]\x00"),
	)
	globalVars.FormatComma = ir.NewGlobalDef(
		"format.comma",
		constant.NewCharArrayFromString(",\x00"),
	)

	libs.GlobalVar = globalVars

	log.DebugMessage("built-in variable declared")
}
