package parser

import (
	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func newNodeParent(kind mTypes.NodeKind, child *mTypes.Node, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind:  kind,
		Child: child,
		Val:   val,
	}
}

func newNodeScalar(ty mTypes.ModoType, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind: mTypes.ND_SCALAR,
		Type: &mTypes.NodeType{
			Value: ty,
		},
		Val: val,
	}
}

func parseExprs(
	rootToken *mTypes.Token,
	exprKind mTypes.NodeKind,
) (*mTypes.Token, *mTypes.Node) {

	var isScalarOrFunc = func(token *mTypes.Token) bool {
		return token.IsKind(mTypes.TK_INT) ||
			token.IsKind(mTypes.TK_FLOAT) ||
			token.IsKind(mTypes.TK_STR) ||
			token.IsKind(mTypes.TK_BOOL) ||
			token.IsKind(mTypes.TK_IDENT) ||
			token.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) ||
			token.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_OPEN)
	}

	nextToken, argHead := parseDeclare(rootToken, exprKind)
	prevNode := argHead
	// NOTE: what means?
	for isScalarOrFunc(nextToken) {
		nextToken, prevNode.Next = parseDeclare(nextToken, exprKind)
		prevNode = prevNode.Next
	}
	return nextToken, argHead
}

func parseBody(
	rootToken *mTypes.Token,
	parentKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {

	nextToken, argHead := parseExprs(rootToken.Next, parentKind)
	// Note: validate argument properties here?
	rootNode := newNodeParent(parentKind, argHead, exprName)
	return nextToken, rootNode
}

func typeCollectionNode(dec *mTypes.Node, tl *mTypes.NodeType) {
	if dec.Kind != mTypes.ND_COLLECTION {
		return
	}

	for d := dec; d != nil; d = d.Next {
		d.Type = tl
		if d.Child != nil {
			typeCollectionNode(d.Child, tl.Child)
		}
	}
}

// TODO: refactoring
func parseIdent(
	tok *mTypes.Token,
	head *mTypes.Node,
	parentKind mTypes.NodeKind,
) (*mTypes.Token, *mTypes.Node) {
	if parentKind == mTypes.ND_DECLARE {
		log.DebugValueColored("is Variable declaration :have %s", tok)

		identName := tok.Val

		type typeListProps struct {
			Type *mTypes.NodeType
			Next *typeListProps
		}
		typeCur := &typeListProps{}
		typeList := typeCur

		// create linked-list(typeList)
		if tok.Next.IsKind(mTypes.TK_TYPE_SIG) {
			tok = tok.Next.Next
			for {
				if !tok.IsKindType() {
					break
				}
				ty, _ := mTypes.GetModoType(tok.Kind)
				typeCur.Type = ty

				if typeCur.Type.Value == mTypes.TY_EXTENDED {
					typeCur.Type.ExtendName = tok.Val
				}

				if tok.Kind.Child != nil {
					typeCur.Type.Child, _ = mTypes.GetModoType(tok.Kind.Child)
				}

				typeCur.Next = &typeListProps{}
				typeCur = typeCur.Next

				tok = tok.Next

				if tok.IsKind(mTypes.TK_TYPE_ARROW) {
					tok = tok.Next
				}
			}
		} else {
			log.Panic("%s : must have '::' with var type", error.ERROR_SYNTAX_ERROR)
		}

		// this parseDeclare which parse all of child expression may return node which have Args
		tok, head = parseDeclare(tok, mTypes.ND_VAR_DECLARE)
		if head.Args != nil {
			for nodeArg := head.Args; nodeArg != nil; nodeArg = nodeArg.Next {
				if typeList.Type == nil {
					log.Panic(
						"%s : mismatch between declared types and function parameters",
						error.ERROR_SYNTAX_ERROR,
					)
				}
				// set Args.Type with correspond linked-list Type
				nodeArg.Type = typeList.Type
				typeList = typeList.Next
			}

		}

		// wrap ND_LAMBDA/ND_VAR_REFERENCE by ND_VAR_DECLARE
		varDeclareNode := newNodeParent(mTypes.ND_VAR_DECLARE, head, identName)
		// HACK: seems buggy
		// the last element of typeList must be return type of function
		varDeclareNode.Type = typeList.Type

		// Only in the case of extended typing that includes structs, the declared type
		// information is treated as the variable’s type as-is.
		// For most other values, especially those involving vectors, the type is not
		// determined by the declared type, but rather inferred from the actual
		// structure—such as the SyntaxRole of the token or its child nodes.
		//
		// Additionally, There are tha plan that the use and definition of struct-typed values are disabled
		// unless they are introduced via a let binding.

		varDecChild := varDeclareNode.Child
		if varDecChild.Type != nil && varDecChild.Type.Value == mTypes.TY_EXTENDED {
			varDecChild.Type.ExtendName = typeList.Type.ExtendName
		}
		return tok, varDeclareNode

	} else {
		log.DebugValueColored("is Variable reference :have %+v", tok)
		return tok.Next, newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)

	}

}

func parseLambda(tok *mTypes.Token, head *mTypes.Node) (*mTypes.Token, *mTypes.Node) {

	// arguments
	argCur := &mTypes.Node{}
	argHead := argCur

	tok = tok.Next

	for tok = tok.Next; !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE); tok = tok.Next {
		if !tok.IsKind(mTypes.TK_IDENT) {
			log.Panic(
				"%s : missing argument symbol, find other type symbol",
				error.ERROR_SYNTAX_ERROR,
			)
		}

		argCur.Next = newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)
		argCur = argCur.Next
	}

	// expressions body
	// NOTE: why passing ]?
	tok, head = parseBody(tok, mTypes.ND_EXPR, "")
	head = newNodeParent(
		mTypes.ND_LAMBDA,
		head,
		"",
	)
	head.Args = argHead.Next

	return tok, head
}

// NOTE: typed at here: ND_SCALAR, ND_VAR_DECLARE, ND_VAR_REFERENCE(Args), ND_EQ, ND_ADD
func parseDeclare(tok *mTypes.Token, parentKind mTypes.NodeKind) (*mTypes.Token, *mTypes.Node) {
	head := &mTypes.Node{}

	if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		tok = tok.Next

		if tok.IsKind(mTypes.TK_DECLARE_VAR) {
			tok, head = parseDeclare(tok.Next, mTypes.ND_DECLARE)
			// A child of ND_DECLARE must be either value(such as string, int, vector) or lambda
			// In case of value, treat it as global scope variable
			if head.Kind == mTypes.ND_VAR_DECLARE && head.Child.Kind != mTypes.ND_LAMBDA {
				head.Child.IsGlobal = true
			}
			head = newNodeParent(mTypes.ND_DECLARE, head, "")

		} else if tok.IsKind(mTypes.TK_DECLARE_TYPE) {
			tok = tok.Next
			head.Kind = mTypes.ND_DECLARE
			structName := tok.Val

			tok = tok.Next
			if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACE_OPEN) {
				log.Panic("%s: missing '{' for defschema", error.ERROR_SYNTAX_ERROR)
			}
			// Note: prefer to be TY_TYPE_STRUCT
			varDeclare := &mTypes.Node{
				Kind: mTypes.ND_TYPE_DECLARE,
				Val:  structName,
				Type: &mTypes.NodeType{Value: mTypes.TY_STRUCT},
			}
			head.Child = varDeclare

			tok = tok.Next
			childHead := &mTypes.Node{}
			child := childHead

			for {
				structElmeName := tok.Val

				tok = tok.Next
				if !tok.IsKind(mTypes.TK_TYPE_SIG) {
					log.Panic("%s: missing type signature '::' for type declaration", error.ERROR_SYNTAX_ERROR)
				}

				tok = tok.Next
				structElemType, _ := mTypes.GetModoType(tok.Kind)
				// Note: Kind needed?
				structTy := &mTypes.Node{
					Val:  structElmeName,
					Type: structElemType,
				}
				child.Next = structTy
				child = child.Next
				tok = tok.Next

				if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACE_CLOSE) {
					tok = tok.Next
					break
				}
				if tok == nil {
					log.Panic("%s: missing closing paren '}'", error.ERROR_SYNTAX_ERROR)
				}
			}
			varDeclare.Child = childHead.Next

		} else if tok.IsKind(mTypes.TK_LAMBDA) {
			tok, head = parseLambda(tok, head)

		} else if tok.IsKind(mTypes.TK_BIND) {
			tok = tok.Next

			if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
				log.Panic("%s: let bindings must involve the more than one pair of variable name and value", error.ERROR_SYNTAX_ERROR)
			}

			prev := &mTypes.Node{}
			varHead := prev

			tok = tok.Next
			for {

				tok, prev.Next = parseDeclare(tok, mTypes.ND_DECLARE)
				prev = prev.Next

				if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
					break
				}

				if tok == nil {
					log.Panic("%s: missing clossing brace ']' for let bindings", error.ERROR_SYNTAX_ERROR)
				}
			}

			nextToken, body := parseBody(tok, mTypes.ND_EXPR, "")
			bind := newNodeParent(mTypes.ND_BIND, body, "")
			bind.Bind = varHead.Next
			// NOTE: why must proceed additionally?
			nextToken = nextToken.Next

			return nextToken, bind

		} else if tok.IsKind(mTypes.TK_LIBCALL) {
			log.DebugValueColored("is Library :have %+v", tok)
			v := tok.Val
			tok, head = parseBody(tok, mTypes.ND_LIBCALL, v)

		} else if tok.IsKind(mTypes.TK_IF) {
			log.DebugValueColored("is IF :have %s", tok)
			head.Kind = mTypes.ND_IF

			nextToken, cond := parseDeclare(tok.Next, mTypes.ND_IF)
			head.Cond = cond

			nextToken, then := parseDeclare(nextToken, mTypes.ND_IF)
			head.Then = newNodeParent(mTypes.ND_EXPR, then, "then")

			nextToken, els := parseDeclare(nextToken, mTypes.ND_IF)
			head.Else = newNodeParent(mTypes.ND_EXPR, els, "els")
			tok = nextToken

		} else if tok.IsKind(mTypes.TK_IDENT) {
			log.DebugValueColored("is calling function :have %+v", tok)
			tok, head = parseBody(tok, mTypes.ND_FUNCCALL, tok.Val)
		}

		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_CLOSE) {
			log.Panic("%s: missing ) for condition expression", error.ERROR_SYNTAX_ERROR)
		}
		return tok.Next, head

	} else if tok.IsKind(mTypes.TK_IDENT) {
		return parseIdent(tok, head, parentKind)

	} else if tok.IsKind(mTypes.TK_INT) {
		return tok.Next, newNodeScalar(mTypes.TY_INT32, tok.Val)

	} else if tok.IsKind(mTypes.TK_FLOAT) {
		return tok.Next, newNodeScalar(mTypes.TY_FLOAT, tok.Val)

	} else if tok.IsKind(mTypes.TK_BOOL) {
		var v string

		switch tok.Val {
		case "true":
			v = "1"
		case "false":
			v = "0"
		default:
			log.Panic("%s: bool value must either 'true' or 'false' : have %s", error.ERROR_SYNTAX_ERROR, tok.Val)
		}
		return tok.Next, newNodeScalar(mTypes.TY_BOOL, v)

		// 'nil' means both of type signature and nil value anyway
	} else if tok.IsKind(mTypes.TK_NIL) {
		return tok.Next, newNodeScalar(mTypes.TY_NIL, "nil")

	} else if tok.IsKind(mTypes.TK_STR) {
		strNode := newNodeScalar(mTypes.TY_STR, tok.Val)
		strNode.Len = uint64(len(tok.Val))
		return tok.Next, strNode

		// means struct value
	} else if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACE_OPEN) {
		tok, rootNode := parseBody(tok, mTypes.ND_COLLECTION, "")
		// ExtendedName (equals to struct type name) is given by parent node,
		rootNode.Type = &mTypes.NodeType{Value: mTypes.TY_EXTENDED}
		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACE_CLOSE) {
			log.Panic("%s: missing close brace for struct", error.ERROR_SYNTAX_ERROR)
		}
		tok = tok.Next
		return tok, rootNode

		// means vector value
	} else if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_OPEN) {
		t, rootNode := parseBody(tok, mTypes.ND_COLLECTION, "")
		rootNode.Type = &mTypes.NodeType{Value: mTypes.TY_VECTOR, Child: rootNode.Child.Type}
		tok = t
		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
			log.Panic("%s: missing close bracket for vector", error.ERROR_SYNTAX_ERROR)
		}
		tok = tok.Next
		return tok, rootNode

	} else {
		log.Panic("%s: unexpected character used, or missing essential signature to parse: have %+v", error.ERROR_SYNTAX_ERROR, tok)
	}

	return tok, head

}

func parseProgram(tok *mTypes.Token) *mTypes.Program {
	p := &mTypes.Program{}
	prevDeclare := p.Declare.Func

	for tok != nil && tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {

		if !tok.Next.IsKind(mTypes.TK_DECLARE_VAR) && !tok.Next.IsKind(mTypes.TK_DECLARE_TYPE) {
			log.Panic(
				"%s: missing declaration symbol, must be either 'def' or 'defschema': have %s",
				error.ERROR_SYNTAX_ERROR,
				tok.Next.Val,
			)
		}

		if prevDeclare == nil {
			tok, p.Declare.Func = parseDeclare(tok, mTypes.ND_PROGRAM_ROOT)
			prevDeclare = p.Declare.Func

		} else {
			tok, prevDeclare.Next = parseDeclare(tok, mTypes.ND_PROGRAM_ROOT)
			prevDeclare = prevDeclare.Next
		}
	}
	return p
}

// take Token object, return Program object
func Parse(token *mTypes.Token) *mTypes.Program {
	log.DebugMessage("code parsing")
	prog := parseProgram(token)
	log.DebugMessage("code parsed")

	prog.Debug(0)

	return prog
}
