package types

import (
	"fmt"
	"strings"
)

var (
	// Scalar
	FRACTIONAL_REG_EXP = `-?\d+\.\d+`
	INTEGER_REG_EXP    = `-?\d+`
	STRING_REG_EXP     = `"([^"]*)"`
	BOOL_REG_EXP       = "\\b(true|false)\\b"

	// Brachets
	PARREN_OPEN      = "("
	PARREN_CLOSE     = ")"
	BRACKET_OPEN     = "["
	BRACKET_CLOSE    = "]"
	BRACE_OPEN       = "{"
	BRACE_CLOSE      = "}"
	BRACKETS_REG_EXP = fmt.Sprintf(
		"[%s%s%s\\%s%s%s]",
		PARREN_OPEN,
		PARREN_CLOSE,
		BRACKET_OPEN,
		BRACKET_CLOSE,
		BRACE_OPEN,
		BRACE_CLOSE,
	)

	// operators
	OPERATOR_ADD = "+"
	OPERATOR_SUB = "-"
	OPERATOR_MUL = "*"
	OPERATOR_DIV = "/"

	OPERATOR_MOD = "mod"

	OPERATOR_EQ = "="
	OPERATOR_LT = "<"
	OPERATOR_GT = ">"

	OPERATOR_AND = "and"
	OPERATOR_OR  = "or"

	OPERATORS_REG_EXP = fmt.Sprintf(
		"[%s\\%s%s%s%s%s%s]|\\b(%s|%s|%s)\\b",
		OPERATOR_ADD,
		OPERATOR_SUB,
		OPERATOR_MUL,
		OPERATOR_DIV,
		OPERATOR_EQ,
		OPERATOR_LT,
		OPERATOR_GT,
		OPERATOR_MOD,
		OPERATOR_AND,
		OPERATOR_OR,
	)

	SYMBOL_THREADING_FIRST = "->"
	SYMBOL_THREADING_LAST  = "->>"
	THREADING_REG_EXP      = fmt.Sprintf("%s|%s", SYMBOL_THREADING_FIRST, SYMBOL_THREADING_LAST)

	SYMBOL_IF      = "if"
	SYMBOL_COND    = "cond"
	BRANCH_REG_EXP = fmt.Sprintf("\\b(%s|%s)\\b", SYMBOL_IF, SYMBOL_COND)

	SYMBOL_DEF       = "def" // NOTE: is core library in clojure
	SYMBOL_DEFSCHEMA = "defschema"
	SYMBOL_FN        = "fn"
	SYMBOL_LET       = "let"
	DECLARE_REG_EXP  = fmt.Sprintf(
		"\\b(%s|%s|%s|%s)\\b",
		SYMBOL_DEF,
		SYMBOL_LET,
		SYMBOL_FN,
		SYMBOL_DEFSCHEMA,
	)

	// Core library
	PRELUDE_PRN              = "prn"
	PRELUDE_GET              = "get"
	PRELUDE_CONJ             = "conj"
	PRELUDE_ASSOC            = "assoc"
	PRELUDE_POP              = "pop"
	PRELUDE_FUNCTION_REG_EXP = fmt.Sprintf(
		"\\b(%s|%s|%s|%s|%s)\\b",
		PRELUDE_PRN,
		PRELUDE_GET,
		PRELUDE_CONJ,
		PRELUDE_ASSOC,
		PRELUDE_POP,
	)

	// Type signature
	SYMBOL_TYPE_SIG   = "::"
	SYMBOL_TYPE_ARROW = "=>"
	SYMBOL_TYPE_INT   = "int"
	SYMBOL_TYPE_STR   = "string"
	SYMBOL_TYPE_NIL   = "nil"
	SYMBOL_TYPE_BOOL  = "bool"

	SYMBOL_TYPE_SCALAR = fmt.Sprintf(
		"(\\b(%s|%s|%s|%s)\\b)",
		SYMBOL_TYPE_INT,
		SYMBOL_TYPE_STR,
		SYMBOL_TYPE_NIL,
		SYMBOL_TYPE_BOOL,
	)
	SYMBOL_TYPE_VECTOR = `\[+(\w+)\]+`

	TYPE_REG_EXP = fmt.Sprintf(
		"(%s|%s|%s|%s)",
		SYMBOL_TYPE_SIG,
		SYMBOL_TYPE_ARROW,
		SYMBOL_TYPE_VECTOR,
		SYMBOL_TYPE_SCALAR,
	)

	SYMBOL_UNDEFINED_REG_EXP = `\w+`

	ALL_REG_EXP = fmt.Sprintf(
		"%s",
		strings.Join(
			[]string{
				FRACTIONAL_REG_EXP,
				INTEGER_REG_EXP,
				STRING_REG_EXP,
				BOOL_REG_EXP,
				TYPE_REG_EXP,
				THREADING_REG_EXP,
				BRANCH_REG_EXP,
				OPERATORS_REG_EXP,
				BRACKETS_REG_EXP,
				PRELUDE_PRN,
				SYMBOL_UNDEFINED_REG_EXP,
			},
			"|",
		),
	)
)
