package parser

import (
	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func validateVarDeclare(n *mTypes.Node, root *mTypes.Node) {
	if n == nil {
		return
	}

	if n.Kind == mTypes.ND_TYPE_DECLARE {
		return
	}

	if n.Kind == mTypes.ND_VAR_DECLARE {
		if n.Child == nil {
			log.Panic(
				"%s: the value of '%s' not defined",
				error.ERROR_SYNTAX_ERROR,
				n.Val,
			)

		} else if n.Child.Type == nil {
			// skip validating
			validateVarDeclare(n.Child, root)

		} else if n.Child.Kind == mTypes.ND_FUNCCALL {
			// skip validating
			validateVarDeclare(n.Child, root)

		} else if n.Type.Value != n.Child.Type.Value {
			log.Panic(
				"%s: cannot use %s (%s type) as %s",
				error.ERROR_SYNTAX_ERROR,
				n.Val,
				n.Type.Value,
				n.Child.Type.Value,
			)
		}
	}
	validateVarDeclare(n.Next, root)
	validateVarDeclare(n.Child, root)
	validateVarDeclare(n.Bind, root)
}

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
