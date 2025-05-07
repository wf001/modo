package parser

import (
	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func validateVarDeclare(n *mTypes.Node, root *mTypes.Node, used map[string]bool) {
	// To check for duplicate variable and type declarations, a Set would be appropriate,
	// but since Go doesn't include one in the standard library,
	// implement provisionally using a map to avoid extra dependencies.
	if n == nil {
		return
	}

	if n.Kind == mTypes.ND_TYPE_DECLARE {
		used[n.Val] = true
		return
	}

	if n.Kind == mTypes.ND_VAR_DECLARE {
		if used[n.Val] {
			log.Panic(
				"%s: cannot use %s, already used",
				error.ERROR_SYNTAX_ERROR,
				n.Val,
			)
		}
		used[n.Val] = true

		if n.Child == nil {
			log.Panic(
				"%s: the value of '%s' not defined",
				error.ERROR_SYNTAX_ERROR,
				n.Val,
			)

		} else if n.Child.Type == nil {
			// skip validating

		} else if n.Child.Kind == mTypes.ND_FUNCCALL {
			// skip validating

		} else if n.Type.Value != n.Child.Type.Value {
			log.Panic(
				"%s: cannot use %s type as %s (%s type)",
				error.ERROR_SYNTAX_ERROR,
				n.Child.Type.Value,
				n.Val,
				n.Type.Value,
			)
		}
	}
	if n.Next != nil {
		validateVarDeclare(n.Next, root, used)
	}
	if n.Child != nil {
		validateVarDeclare(n.Child, root, used)
	}
	if n.Bind != nil {
		validateVarDeclare(n.Bind, root, used)
	}
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
			log.Panic("%s: undefined: %s", error.ERROR_SYNTAX_ERROR, n.Val)
		}
	}
	validateReference(n.Next, root)
	validateReference(n.Child, root)
	validateReference(n.Bind, root)
}

func findTypeDeclare(n *mTypes.Node, targetType string) *mTypes.Node {
	if n == nil {
		return nil
	}
	if n.Kind == mTypes.ND_TYPE_DECLARE && n.Val == targetType {
		return n
	}
	if result := findTypeDeclare(n.Child, targetType); result != nil {
		return result
	}
	return nil
}

func validateExtendedTypeReference(n *mTypes.Node, root *mTypes.Node) {
	if n == nil {
		return
	}
	if n.Kind == mTypes.ND_VAR_DECLARE && n.Type.Value == mTypes.TY_EXTENDED {
		decl := findTypeDeclare(root, n.Type.ExtendName)
		if decl == nil {
			log.Panic("%s: undefined type: %s", error.ERROR_SYNTAX_ERROR, n.Type.ExtendName)
		}
	}
	validateExtendedTypeReference(n.Next, root)
	validateExtendedTypeReference(n.Child, root)
	validateExtendedTypeReference(n.Bind, root)
}
