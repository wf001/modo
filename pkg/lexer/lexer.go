package lexer

import (
	"regexp"
	"strings"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

// defines the structure for token patterns as a linked list
type tokenPattern struct {
	Pattern   string
	TokenType string
	Next      *tokenPattern
}

func isMatched(s string, typ string) bool {
	re := regexp.MustCompile(typ)
	return re.MatchString(s)
}

// traverses the linked list to find a matching pattern and returns the TokenType
func (tp *tokenPattern) matchTokenType(s string) (string, bool) {
	current := tp
	for current != nil {
		if isMatched(s, current.Pattern) {
			return current.TokenType, true
		}
		current = current.Next
	}
	return "", false
}

// initializes a TokenPattern with predefined patterns and token types
func newTokenPattern() *tokenPattern {
	head := &tokenPattern{
		Pattern:   mTypes.INTEGER_REG_EXP,
		TokenType: mTypes.TK_INT,
	}
	current := head

	add := func(pattern, tokenType string) {
		newNode := &tokenPattern{
			Pattern:   pattern,
			TokenType: tokenType,
		}
		current.Next = newNode
		current = newNode
	}

	add(mTypes.SYMBOL_TYPE_SIG, mTypes.TK_TYPE_SIG)
	add(mTypes.SYMBOL_TYPE_ARROW, mTypes.TK_TYPE_ARROW)
	add(mTypes.SYMBOL_TYPE_VECTOR, mTypes.TK_TYPE_VECTOR)
	add(mTypes.SYMBOL_TYPE_INT, mTypes.TK_TYPE_INT)
	add(mTypes.SYMBOL_TYPE_STR, mTypes.TK_TYPE_STR)
	add(mTypes.SYMBOL_TYPE_NIL, mTypes.TK_TYPE_NIL)
	add(mTypes.SYMBOL_TYPE_BOOL, mTypes.TK_TYPE_BOOL)

	add(mTypes.STRING_REG_EXP, mTypes.TK_STR)
	add(mTypes.BOOL_REG_EXP, mTypes.TK_BOOL)
	add(mTypes.OPERATORS_REG_EXP, mTypes.TK_LIBCALL)
	add(mTypes.BRACKETS_REG_EXP, mTypes.TK_PAREN)
	add(mTypes.PRELUDE_FUNCTION_REG_EXP, mTypes.TK_LIBCALL)
	add(mTypes.SYMBOL_FN, mTypes.TK_LAMBDA)
	add(mTypes.SYMBOL_DEFSCHEMA, mTypes.TK_DECLARE_TYPE)
	add(mTypes.SYMBOL_DEF, mTypes.TK_DECLARE_VAR)
	add(mTypes.SYMBOL_LET, mTypes.TK_BIND)
	add(mTypes.SYMBOL_IF, mTypes.TK_IF)

	return head
}

func newToken(kind mTypes.TokenKind, prev *mTypes.Token, val string) *mTypes.Token {
	tok := &mTypes.Token{
		Kind: kind,
		Val:  val,
	}
	prev.Next = tok
	return tok
}

// remove escaped \" from origin string
func trimQuote(head *mTypes.Token) {
	for t := head.Next; t != nil; t = t.Next {
		if t.IsKind(mTypes.TK_STR) {
			s := strings.Trim(t.Val, "\"")
			s = s + string([]byte{0})
			t.Val = s
		}
	}
}

// In modo, 'nil' denotes both a type name and a value of type nil, all 'nil' are interpreted as type signature at first.
// so 'nil' which is used in a non-type declaration context is reinterpreted as the value nil with this function.
func accurateNilType(head *mTypes.Token) {
	for t := head.Next; t.Next != nil; t = t.Next {
		if !(t.IsKind(mTypes.TK_TYPE_ARROW) || t.IsKind(mTypes.TK_TYPE_SIG)) &&
			t.Next.IsKind(mTypes.TK_TYPE_NIL) {
			t.Next.Kind = mTypes.TK_NIL
		}
	}
}

func splitString(expr string) []string {
	re := regexp.MustCompile(mTypes.STRING_REG_EXP)
	expr = re.ReplaceAllString(expr, `"$1"`)

	log.Debug(log.YELLOW("preprocessed %#+v"), expr)

	re = regexp.MustCompile(mTypes.ALL_REG_EXP)
	log.Debug(log.YELLOW("regex to analyse: %#+v"), mTypes.ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted program: %#+v"), res)
	return res
}

// performs lexical analysis on the input strings
func doLexicalAnalyse(splittedString []string) *mTypes.Token {
	prev := &mTypes.Token{}
	head := prev

	tokenMap := newTokenPattern()

	for _, p := range splittedString {
		if tokenType, matched := tokenMap.matchTokenType(p); matched {
			prev = newToken(tokenType, prev, p)
			// It seems that Go's regexp package does not support lookbehind or lookahead assertions like "(?<!...)", "(?=...)",
			// so it is difficult to write a regex that excludes numbers that are
			// immediately preceded by letters (e.g., to avoid matching the "2" in "vec2").
			//
			// Instead, previous processing (splittedString) match numbers broadly using patterns like `-?\d+`,
			// and rely on post-processing logic to correctly classify tokens
			if prev.IsKind(mTypes.TK_INT) {
				numRe := regexp.MustCompile(`^-?\d+$`)

				isNumber := numRe.MatchString(p)
				if !isNumber {
					log.Debug(log.YELLOW("change token kind to TY_IDENT: %#+v"), prev)
					prev.Kind = mTypes.TK_IDENT
				}

			}
			// when token is vector, set element type of the vector
			if prev.IsKind(mTypes.TK_TYPE_VECTOR) {
				vecRe := regexp.MustCompile(`\[(\w+)\]`)
				matches := vecRe.FindAllStringSubmatch(prev.Val, -1)

				for _, match := range matches {
					if tokenType, matched := tokenMap.matchTokenType(match[1]); matched {
						prev.ChildKind = tokenType
					}
				}
			}
		} else {
			log.Debug("regard '%+v' as variable declaration or reference symbol", p)
			prev = newToken(mTypes.TK_IDENT, prev, p)
		}
	}
	trimQuote(head)
	accurateNilType(head)

	head = head.Next
	head.DebugTokens()
	return head
}

// take string, return Token object
func Lex(s string) *mTypes.Token {
	log.Debug(log.YELLOW("original source: '%s'"), s)
	log.DebugMessage("code lexing")

	arr := splitString(s)
	res := doLexicalAnalyse(arr)

	log.DebugMessage("code lexed")
	return res
}
