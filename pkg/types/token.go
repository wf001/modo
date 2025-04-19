package types

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
)

type TokenKind = string

const (
	TK_DECLARE_VAR  = TokenKind("TK_DECLARE_VAR")
	TK_DECLARE_TYPE = TokenKind("TK_DECLARE_TYPE")
	TK_LAMBDA       = TokenKind("TK_LAMBDA")
	TK_BIND         = TokenKind("TK_BIND")
	TK_IF           = TokenKind("TK_IF")

	TK_PAREN = TokenKind("TK_PAREN")

	TK_LIBCALL = TokenKind("TK_LIBCALL")
	TK_IDENT   = TokenKind("TK_IDENT")

	TK_TYPE_SIG      = TokenKind("TK_TYPE_SIG")
	TK_TYPE_ARROW    = TokenKind("TK_TYPE_ARROW")
	TK_TYPE_INT      = TokenKind("TK_TYPE_INT")
	TK_TYPE_FLOAT    = TokenKind("TK_TYPE_FLOAT")
	TK_TYPE_STR      = TokenKind("TK_TYPE_STR")
	TK_TYPE_BOOL     = TokenKind("TK_TYPE_BOOL")
	TK_TYPE_NIL      = TokenKind("TK_TYPE_NIL")
	TK_TYPE_VECTOR   = TokenKind("TK_TYPE_VECTOR")
	TK_TYPE_STRUCT   = TokenKind("TK_TYPE_STRUCT")
	TK_TYPE_EXTENDED = TokenKind("TK_TYPE_EXTENDED")

	TK_INT   = TokenKind("TK_INT")
	TK_FLOAT = TokenKind("TK_FLOAT")
	TK_BOOL  = TokenKind("TK_BOOL")
	TK_STR   = TokenKind("TK_STR")
	TK_NIL   = TokenKind("TK_NIL")
)

type Token struct {
	Kind      TokenKind
	ChildKind TokenKind // NOTE: use a different Kind Struct?
	Next      *Token
	Val       string
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"{Val:\"%s\", Kind:\"%s\", ChildKind:\"%s\"}",
		t.Val,
		t.Kind,
		t.ChildKind,
	)
}

// ==============
// predication
// ==============

func (tok *Token) IsKindAndVal(kind string, val string) bool {
	return tok != nil && tok.IsKind(kind) && tok.Val == val
}

func (tok *Token) IsKind(kind TokenKind) bool {
	return tok.Kind == kind
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

func GetModoType(k TokenKind) (ModoType, bool) {
	var typeMap = map[string]ModoType{
		TK_TYPE_INT:      TY_INT32,
		TK_TYPE_STR:      TY_STR,
		TK_TYPE_NIL:      TY_NIL,
		TK_TYPE_BOOL:     TY_BOOL,
		TK_TYPE_VECTOR:   TY_VECTOR,
		TK_TYPE_STRUCT:   TY_STRUCT,
		TK_TYPE_EXTENDED: TY_EXTENDED,
	}

	if kind, exists := typeMap[k]; exists {
		return kind, true
	}

	return "", false
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
