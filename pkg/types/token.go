package types

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
)

type SyntaxRole = string

const (
	TK_DECLARE_VAR  = SyntaxRole("TK_DECLARE_VAR")
	TK_DECLARE_TYPE = SyntaxRole("TK_DECLARE_TYPE")
	TK_LAMBDA       = SyntaxRole("TK_LAMBDA")
	TK_BIND         = SyntaxRole("TK_BIND")
	TK_IF           = SyntaxRole("TK_IF")

	TK_PAREN = SyntaxRole("TK_PAREN")

	TK_LIBCALL = SyntaxRole("TK_LIBCALL")
	TK_IDENT   = SyntaxRole("TK_IDENT")

	TK_TYPE_SIG      = SyntaxRole("TK_TYPE_SIG")
	TK_TYPE_ARROW    = SyntaxRole("TK_TYPE_ARROW")
	TK_TYPE_INT      = SyntaxRole("TK_TYPE_INT")
	TK_TYPE_FLOAT    = SyntaxRole("TK_TYPE_FLOAT")
	TK_TYPE_STR      = SyntaxRole("TK_TYPE_STR")
	TK_TYPE_BOOL     = SyntaxRole("TK_TYPE_BOOL")
	TK_TYPE_NIL      = SyntaxRole("TK_TYPE_NIL")
	TK_TYPE_VECTOR   = SyntaxRole("TK_TYPE_VECTOR")
	TK_TYPE_STRUCT   = SyntaxRole("TK_TYPE_STRUCT")
	TK_TYPE_EXTENDED = SyntaxRole("TK_TYPE_EXTENDED")

	TK_INT   = SyntaxRole("TK_INT")
	TK_FLOAT = SyntaxRole("TK_FLOAT")
	TK_BOOL  = SyntaxRole("TK_BOOL")
	TK_STR   = SyntaxRole("TK_STR")
	TK_NIL   = SyntaxRole("TK_NIL")
)

type TokenKind struct {
	Value SyntaxRole
	Child *TokenKind
}

func (t *TokenKind) String() string {
	if t.Child != nil {
		return fmt.Sprintf(
			"{Value:\"%s\", Child: \"%s\" }",
			t.Value,
			t.Child.String(),
		)

	}
	return fmt.Sprintf(
		"{Value:\"%s\"}",
		t.Value,
	)
}

type Token struct {
	Kind *TokenKind
	Next *Token
	Val  string
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"{Val:\"%s\", Kind:\"%s\" }",
		t.Val,
		t.Kind,
	)
}

// ==============
// predication
// ==============

func (tok *Token) IsKindAndVal(kind string, val string) bool {
	return tok != nil && tok.IsKind(kind) && tok.Val == val
}

func (tok *Token) IsKind(kind SyntaxRole) bool {
	return tok.Kind.Value == kind
}

func (tok *Token) IsKindType() bool {
	return tok.IsKind(TK_TYPE_ARROW) ||
		tok.IsKind(TK_TYPE_INT) ||
		tok.IsKind(TK_TYPE_STR) ||
		tok.IsKind(TK_TYPE_NIL) ||
		tok.IsKind(TK_TYPE_BOOL) ||
		tok.IsKind(TK_TYPE_VECTOR) ||
		tok.IsKind(TK_TYPE_EXTENDED)
}

// ==============
// conversion llir/llvm properties to other properties
// ==============

func GetModoType(tkkind *TokenKind) (*NodeType, bool) {
	var scalarType = map[string]ModoType{
		TK_TYPE_INT:      TY_INT32,
		TK_TYPE_STR:      TY_STR,
		TK_TYPE_NIL:      TY_NIL,
		TK_TYPE_BOOL:     TY_BOOL,
		TK_TYPE_EXTENDED: TY_EXTENDED,
	}
	var collType = map[string]ModoType{
		TK_TYPE_VECTOR: TY_VECTOR,
		TK_TYPE_STRUCT: TY_STRUCT,
	}

	if kind, isScalar := scalarType[tkkind.Value]; isScalar {
		return &NodeType{Value: kind}, true
	}

	if parentKind, isColl := collType[tkkind.Value]; isColl {
		childType, _ := GetModoType(tkkind.Child)
		return &NodeType{Value: parentKind, Child: childType}, true
	}

	return &NodeType{Value: ""}, false
}

// ==============
// The following is for developement purposes
// ==============

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.DebugNoLine(log.BLUE(fmt.Sprintf("\t %s", tok)))
	}
}
