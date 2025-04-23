package types

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
)

type SyntaxRole = string

const (
	TK_DECLARE_VAR  = SyntaxRole("tk.declare.var")
	TK_DECLARE_TYPE = SyntaxRole("tk.declare.type")
	TK_LAMBDA       = SyntaxRole("tk.declare.lambda")
	TK_BIND         = SyntaxRole("tk.bind")
	TK_IF           = SyntaxRole("tk.if")

	TK_PAREN = SyntaxRole("tk.paren")

	TK_LIBCALL = SyntaxRole("tk.libcall")
	TK_IDENT   = SyntaxRole("tk.ident")

	TK_TYPE_SIG      = SyntaxRole("tk.type.sig")
	TK_TYPE_ARROW    = SyntaxRole("tk.type.arrow")
	TK_TYPE_INT      = SyntaxRole("tk.type.int")
	TK_TYPE_FLOAT    = SyntaxRole("tk.type.float")
	TK_TYPE_STR      = SyntaxRole("tk.type.str")
	TK_TYPE_BOOL     = SyntaxRole("tk.type.bool")
	TK_TYPE_NIL      = SyntaxRole("tk.type.nil")
	TK_TYPE_VECTOR   = SyntaxRole("tk.type.vector")
	TK_TYPE_STRUCT   = SyntaxRole("tk.type.struct")
	TK_TYPE_EXTENDED = SyntaxRole("tk.type.extended")

	TK_INT   = SyntaxRole("tk.int")
	TK_FLOAT = SyntaxRole("tk.float")
	TK_BOOL  = SyntaxRole("tk.bool")
	TK_STR   = SyntaxRole("tk.str")
	TK_NIL   = SyntaxRole("tk.nil")
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
