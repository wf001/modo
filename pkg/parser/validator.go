package parser

import (
	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func findDeclare(n *mTypes.Node, targetVal string) *mTypes.Node {
	if n == nil {
		return nil
	}
	if n.Kind == mTypes.ND_VAR_DECLARE && n.Val == targetVal {
		return n
	}
	if result := findDeclare(n.Next, targetVal); result != nil {
		return result
	}
	if result := findDeclare(n.Child, targetVal); result != nil {
		return result
	}
	if result := findDeclare(n.Bind, targetVal); result != nil {
		return result
	}
	if result := findDeclare(n.Args, targetVal); result != nil {
		return result
	}
	return nil
}

func validateReference(n *mTypes.Node, root *mTypes.Node) {
	if n == nil {
		return
	}
	if n.Kind == mTypes.ND_VAR_REFERENCE {
		decl := findDeclare(root, n.Val)
		if decl == nil {
			log.Panic("%s: undefined variable: %s", error.ERROR_SYNTAX_ERROR, n.Val)
		}
	}
	validateReference(n.Next, root)
	validateReference(n.Child, root)
	validateReference(n.Bind, root)
}
