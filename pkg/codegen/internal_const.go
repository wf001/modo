package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareInternalConst(ir *ir.Module, internal *mTypes.Internal) {
	globalConst := &mTypes.GlobalConst{}

	globalConst.FormatDigit = ir.NewGlobalDef(
		"_format.digit",
		constant.NewCharArrayFromString("%d\x00"),
	)
	globalConst.FormatStr = ir.NewGlobalDef(
		"_format.string",
		constant.NewCharArrayFromString("%s\x00"),
	)
	globalConst.StringSpace = ir.NewGlobalDef(
		"_constant.string.space",
		constant.NewCharArrayFromString(" \x00"),
	)
	globalConst.StringCR = ir.NewGlobalDef(
		"_constant.string.cr",
		constant.NewCharArrayFromString("\n\x00"),
	)
	globalConst.StringTrue = ir.NewGlobalDef(
		"_constant.string.true",
		constant.NewCharArrayFromString("true\x00"),
	)
	globalConst.StringFalse = ir.NewGlobalDef(
		"_constant.string.false",
		constant.NewCharArrayFromString("false\x00"),
	)
	globalConst.StringNil = ir.NewGlobalDef(
		"_constant.string.nil",
		constant.NewCharArrayFromString("nil\x00"),
	)
	globalConst.StringBracketOpen = ir.NewGlobalDef(
		"_constant.string.bracket.open",
		constant.NewCharArrayFromString("[\x00"),
	)
	globalConst.StringBracketClose = ir.NewGlobalDef(
		"_constant.string.bracket.close",
		constant.NewCharArrayFromString("]\x00"),
	)
	globalConst.StringComma = ir.NewGlobalDef(
		"_constant.string.comma",
		constant.NewCharArrayFromString(",\x00"),
	)

	internal.GlobalConst = globalConst

	log.DebugMessage("built-in variable declared")
}
