package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mTypes "github.com/wf001/modo/pkg/types"
)

func TestSplitProgram(t *testing.T) {
	assert.ElementsMatch(t, []string{}, splitString(""))
	assert.ElementsMatch(t, []string{"1"}, splitString("1"))
	assert.ElementsMatch(t, []string{"+", "123", "45"}, splitString("+ 123 45"))
	assert.ElementsMatch(t, []string{"-", "123", "45"}, splitString("- 123 45"))
	assert.ElementsMatch(t, []string{"*", "123", "45"}, splitString("* 123 45"))
	assert.ElementsMatch(t, []string{"/", "123", "45"}, splitString("/ 123 45"))
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "123", "45", ")"},
		splitString("(+ 123 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "a", "45", ")"},
		splitString("(+ a 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "ae", "45", ")"},
		splitString("(+ ae 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "ifa", "45", ")"},
		splitString("(+ ifa 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "->", "(", "+", "123", "45", ")", ")"},
		splitString("(->(+ 123 45))"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "=", "1", "1", ")",
		},
		splitString("(= 1 1)"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "def", "main", "(", "fn", "[", "]", "(", "if", "(", "=", "1", "1", ")", "(", "prn",
			"(", "+", "1", "2", ")", ")", "(", "prn", "(", "+", "3", "-3", ")", ")", ")", ")", ")",
		},
		splitString("(def main (fn [] (if (= 1 1) (prn (+ 1 2)) (prn (+ 3 -3)))))"),
	)
	t.Log(mTypes.ALL_REG_EXP)
}

func TestLexOneInteger(t *testing.T) {
	assert.Equal(
		t,
		&mTypes.Token{Kind: &mTypes.TokenKind{Value: mTypes.TK_INT}, Val: "1"},
		Lex("1"),
	)
}

func add(kind mTypes.SyntaxRole, val string) *mTypes.Token {
	return &mTypes.Token{
		Kind: &mTypes.TokenKind{Value: kind},
		Val:  val,
	}
}

func buildToken(tokens []*mTypes.Token) *mTypes.Token {
	for i := 0; i < len(tokens)-1; i++ {
		tokens[i].Next = tokens[i+1]
	}
	return tokens[0]
}

func TestLexOperationAdd(t *testing.T) {
	tokens := []*mTypes.Token{
		add(mTypes.TK_PAREN, "("),
		add(mTypes.TK_LIBCALL, "+"),
		add(mTypes.TK_INT, "1"),
		add(mTypes.TK_INT, "2"),
		add(mTypes.TK_PAREN, ")"),
	}

	want := buildToken(tokens)

	assert.EqualValues(t, want, Lex("(+ 1 2)"))
}

func TestLexOperationAddTakingAdd(t *testing.T) {
	tokens := []*mTypes.Token{
		add(mTypes.TK_PAREN, "("),
		add(mTypes.TK_LIBCALL, "+"),
		add(mTypes.TK_PAREN, "("),
		add(mTypes.TK_LIBCALL, "+"),
		add(mTypes.TK_INT, "1"),
		add(mTypes.TK_INT, "2"),
		add(mTypes.TK_PAREN, ")"),
		add(mTypes.TK_PAREN, "("),
		add(mTypes.TK_LIBCALL, "+"),
		add(mTypes.TK_INT, "3"),
		add(mTypes.TK_INT, "4"),
		add(mTypes.TK_PAREN, ")"),
		add(mTypes.TK_PAREN, ")"),
	}

	want := buildToken(tokens)

	assert.EqualValues(t, want, Lex("(+ (+ 1 2) (+ 3 4))"))
}

func TestNewTokenMap(t *testing.T) {
	res := newTokenPattern()
	assert.EqualValues(t, mTypes.INTEGER_REG_EXP, res.Pattern)
	assert.EqualValues(t, mTypes.SYMBOL_TYPE_SIG, res.Next.Pattern)
}
