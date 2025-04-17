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
	globalConst.FormatDigit.Immutable = true

	globalConst.FormatStr = ir.NewGlobalDef(
		"_format.string",
		constant.NewCharArrayFromString("%s\x00"),
	)
	globalConst.FormatStr.Immutable = true

	globalConst.StringSpace = ir.NewGlobalDef(
		"_constant.string.space",
		constant.NewCharArrayFromString(" \x00"),
	)
	globalConst.StringSpace.Immutable = true

	globalConst.StringCR = ir.NewGlobalDef(
		"_constant.string.cr",
		constant.NewCharArrayFromString("\n\x00"),
	)
	globalConst.StringCR.Immutable = true

	globalConst.StringTrue = ir.NewGlobalDef(
		"_constant.string.true",
		constant.NewCharArrayFromString("true\x00"),
	)
	globalConst.StringTrue.Immutable = true

	globalConst.StringFalse = ir.NewGlobalDef(
		"_constant.string.false",
		constant.NewCharArrayFromString("false\x00"),
	)
	globalConst.StringFalse.Immutable = true

	globalConst.StringNil = ir.NewGlobalDef(
		"_constant.string.nil",
		constant.NewCharArrayFromString("nil\x00"),
	)
	globalConst.StringNil.Immutable = true

	globalConst.StringBracketOpen = ir.NewGlobalDef(
		"_constant.string.bracket.open",
		constant.NewCharArrayFromString("[\x00"),
	)
	globalConst.StringBracketOpen.Immutable = true

	globalConst.StringBracketClose = ir.NewGlobalDef(
		"_constant.string.bracket.close",
		constant.NewCharArrayFromString("]\x00"),
	)
	globalConst.StringBracketClose.Immutable = true

	globalConst.StringComma = ir.NewGlobalDef(
		"_constant.string.comma",
		constant.NewCharArrayFromString(",\x00"),
	)
	globalConst.StringComma.Immutable = true

	internal.GlobalConst = globalConst

	log.DebugMessage("built-in variable declared")
}
