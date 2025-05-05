package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModoType(t *testing.T) {
	// int32
	tkkind := &TokenKind{Value: TK_TYPE_INT}
	want := &NodeType{Value: TY_INT32}
	res, ok := GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.Equal(t, want, res)

	// bool
	tkkind = &TokenKind{Value: TK_TYPE_BOOL}
	want = &NodeType{Value: TY_BOOL}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.Equal(t, want, res)

	// string
	tkkind = &TokenKind{Value: TK_TYPE_STR}
	want = &NodeType{Value: TY_STR}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.Equal(t, want, res)

	// vector int
	tkkind = &TokenKind{Value: TK_TYPE_VECTOR, Child: &TokenKind{Value: TK_TYPE_INT}}
	want = &NodeType{Value: TY_VECTOR, Child: &NodeType{Value: TY_INT32}}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.EqualValues(t, want, res)

	// vector string
	tkkind = &TokenKind{Value: TK_TYPE_VECTOR, Child: &TokenKind{Value: TK_TYPE_STR}}
	want = &NodeType{Value: TY_VECTOR, Child: &NodeType{Value: TY_STR}}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.EqualValues(t, want, res)

	// vector vector int
	tkkind = &TokenKind{
		Value: TK_TYPE_VECTOR,
		Child: &TokenKind{Value: TK_TYPE_VECTOR, Child: &TokenKind{Value: TK_TYPE_INT}},
	}
	want = &NodeType{
		Value: TY_VECTOR,
		Child: &NodeType{Value: TY_VECTOR, Child: &NodeType{Value: TY_INT32}},
	}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.EqualValues(t, want, res)

	// struct
	tkkind = &TokenKind{Value: TK_TYPE_EXTENDED}
	want = &NodeType{Value: TY_EXTENDED}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, true, ok)
	assert.EqualValues(t, want, res)

	// struct
	tkkind = &TokenKind{Value: ""}
	want = &NodeType{Value: ""}
	res, ok = GetModoType(tkkind)
	assert.Equal(t, false, ok)
	assert.EqualValues(t, want, res)
}
