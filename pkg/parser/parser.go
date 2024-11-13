package parser

import (
	"fmt"

	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/types"
)

func newNode(kind types.NodeKind, child *types.Node) *types.Node {
	return &types.Node{
		Kind:  kind,
		Child: child,
	}
}

func newNodeNum(tok *types.Token) *types.Node {
	return &types.Node{
		Kind: types.ND_INT,
		Val:  tok.Val,
	}
}

func expr(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &types.Node{}

	if tok.Kind == types.TK_EOL {
		return tok, nil
	}
	if tok.Val == "(" {
		tok = tok.Next
		if tok.Val == "+" {
			nextToken, childHead := expr(tok.Next)
			prevNode := childHead
			for nextToken.Kind == types.TK_NUM || nextToken.Val == "(" {
				nextToken, prevNode.Next = expr(nextToken)
				prevNode = prevNode.Next
			}
			tok = nextToken
			head = newNode(types.ND_ADD, childHead)
		}
		if tok.Val != ")" {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head
	} else if tok.Kind == types.TK_NUM {
		return tok.Next, newNodeNum(tok)
	}
	return tok, head
}

// return Node object from Token array
func Parse(tok *types.Token) *types.Node {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	_, node := expr(tok)
	log.DebugNode(node, 0)

	return node
}
