package lexer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTokenizer struct {
	tokens        []Token
	optionalError error
}

func (m *mockTokenizer) Token() (Token, error) {
	if m.optionalError != nil {
		return Token{}, m.optionalError
	}

	if len(m.tokens) == 0 {
		return Token{}, errors.New("no tokens")
	}
	next := m.tokens[0]
	m.tokens = m.tokens[1:]
	return next, nil
}

func Test_peekable_Peek_first(t *testing.T) {
	firstToken := Token{Type: TokenTypeColon}
	mock := &mockTokenizer{tokens: []Token{
		firstToken,
		{Type: TokenTypeComma},
		{Type: TokenTypeEOF},
	}}

	p := &peekable{tokenizer: mock}

	token, err := p.Peek()
	assert.NoError(t, err)
	assert.Equal(t, firstToken, token)
}

func Test_peekable_Peek_mutliple_times(t *testing.T) {
	mock := &mockTokenizer{tokens: []Token{
		{Type: TokenTypeComma},
		{Type: TokenTypeEOF},
	}}

	p := &peekable{tokenizer: mock}

	peeked, err := p.Peek()
	assert.NoError(t, err)
	for i := 0; i < 7; i++ {
		token, err := p.Peek()
		assert.NoError(t, err)
		assert.Equal(t, peeked, token)
	}
}

func Test_peekable_Peek_then_take(t *testing.T) {
	firstToken := Token{Type: TokenTypeColon}
	mock := &mockTokenizer{tokens: []Token{
		firstToken,
		{Type: TokenTypeComma},
		{Type: TokenTypeEOF},
	}}

	p := &peekable{tokenizer: mock}

	peeked, err := p.Peek()
	assert.NoError(t, err)

	// should also return same token
	next, err := p.Token()
	assert.NoError(t, err)

	assert.Equal(t, next, peeked)
}

func Test_peekable_Peek_first_take(t *testing.T) {
	secondToken := Token{Type: TokenTypeColon}
	mock := &mockTokenizer{tokens: []Token{
		{Type: TokenTypeComma},
		secondToken,
		{Type: TokenTypeEOF},
	}}

	p := &peekable{tokenizer: mock}

	_, err := p.Token()
	assert.NoError(t, err)

	peeked, err := p.Peek()
	assert.NoError(t, err)

	assert.Equal(t, secondToken, peeked)
}

func Test_peekable_Peek_error(t *testing.T) {
	expectedErr := errors.New("fail")
	mock := &mockTokenizer{optionalError: expectedErr}

	p := &peekable{tokenizer: mock}

	_, err1 := p.Token()
	assert.Error(t, err1)

	_, err2 := p.Peek()
	assert.Equal(t, err1, err2)
}

func Test_peekable_Token_error(t *testing.T) {
	expectedErr := errors.New("fail")
	mock := &mockTokenizer{optionalError: expectedErr}

	p := &peekable{tokenizer: mock}

	_, err1 := p.Token()
	assert.Error(t, err1)

	_, err2 := p.Token()
	assert.Equal(t, err1, err2)
}
