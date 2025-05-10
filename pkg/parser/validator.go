package parser

import (
	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func findVarDeclare(n *mTypes.Node, targetVal string) *mTypes.Node {
	if n == nil {
		return nil
	}
	if n.IsKind(mTypes.ND_VAR_DECLARE) && n.Val == targetVal {
		return n
	}
	if result := findVarDeclare(n.Next, targetVal); result != nil {
		return result
	}
	if result := findVarDeclare(n.Child, targetVal); result != nil {
		return result
	}
	if result := findVarDeclare(n.Bind, targetVal); result != nil {
		return result
	}
	if result := findVarDeclare(n.Args, targetVal); result != nil {
		return result
	}
	return nil
}

func findTypeDeclare(n *mTypes.Node, targetType string) *mTypes.Node {
	if n == nil {
		return nil
	}
	if n.IsKind(mTypes.ND_TYPE_DECLARE) && n.Val == targetType {
		return n
	}
	if result := findTypeDeclare(n.Next, targetType); result != nil {
		return result
	}
	if result := findTypeDeclare(n.Child, targetType); result != nil {
		return result
	}
	return nil
}

func walkNode(
	n *mTypes.Node,
	root *mTypes.Node,
	used map[string]bool,
	f func(n *mTypes.Node, root *mTypes.Node, used map[string]bool) bool,
) {
	if n == nil {
		return
	}

	if !f(n, root, used) {
		return
	}

	walkNode(n.Next, root, used, f)
	walkNode(n.Child, root, used, f)
	walkNode(n.Bind, root, used, f)
}

func validateNode(root *mTypes.Node) {
	// To check for duplicate variable and type declarations, a Set would be appropriate,
	// but since Go doesn't include one in the standard library,
	// implement provisionally using a map to avoid extra dependencies.

	// checks that all variable references point to valid declarations.
	walkNode(root, root, map[string]bool{},
		func(n *mTypes.Node, root *mTypes.Node, used map[string]bool) bool {
			if n.IsKind(mTypes.ND_VAR_REFERENCE) {
				decl := findVarDeclare(root, n.Val)
				if decl == nil {
					log.Panic("%s: undefined: %s", error.ERROR_SYNTAX_ERROR, n.Val)
				}
			}
			return true
		},
	)
	// checks for duplicate variable/type declarations and type mismatches.
	walkNode(root, root, map[string]bool{},
		func(n *mTypes.Node, root *mTypes.Node, used map[string]bool) bool {

			if n.IsKind(mTypes.ND_TYPE_DECLARE) {
				used[n.Val] = true
				return false
			}

			if n.IsKind(mTypes.ND_VAR_DECLARE) {
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

				} else if n.Child.IsKind(mTypes.ND_FUNCCALL) {
					// skip validating

				} else if !n.IsType(n.Child.Type.Value) {
					log.Panic(
						"%s: cannot use %s type as %s (%s type)",
						error.ERROR_SYNTAX_ERROR,
						n.Child.Type.Value,
						n.Val,
						n.Type.Value,
					)
				}
			}
			return true
		},
	)

	// ensures that extended type references are defined.
	walkNode(root, root, map[string]bool{},
		func(n *mTypes.Node, root *mTypes.Node, used map[string]bool) bool {
			if n.IsKind(mTypes.ND_VAR_DECLARE) && n.IsType(mTypes.TY_EXTENDED) {
				decl := findTypeDeclare(root, n.Type.ExtendName)
				if decl == nil {
					log.Panic("%s: undefined type: %s", error.ERROR_SYNTAX_ERROR, n.Type.ExtendName)
				}
			}
			return true
		},
	)

}
