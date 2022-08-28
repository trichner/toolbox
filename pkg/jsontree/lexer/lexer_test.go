package lexer

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Token_SkipWhitespace(t *testing.T) {
	raw := "   \t "
	lex := &lexer{r: bufio.NewReader(strings.NewReader(raw))}

	_, err := lex.Token()
	assert.NoError(t, err)
}

func TestLexer_Token_WithSuccess(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []Token
	}{
		{
			name:     "emtpy",
			raw:      "",
			expected: []Token{{Type: TokenTypeEOF}},
		},
		{
			name:     "whitespace only",
			raw:      "  ",
			expected: []Token{{Type: TokenTypeEOF}},
		},
		{
			name:     "whitespace and primitive text 1",
			raw:      "  hello  ",
			expected: []Token{{Type: TokenTypePrimitiveText, Value: "hello"}, {Type: TokenTypeEOF}},
		},
		{
			name:     "whitespace and primitive text 2",
			raw:      "  true",
			expected: []Token{{Type: TokenTypePrimitiveText, Value: "true"}, {Type: TokenTypeEOF}},
		},
		{
			name: "array 1",
			raw:  "[true]",
			expected: []Token{
				{Type: TokenTypeOpeningBracket},
				{Type: TokenTypePrimitiveText, Value: "true"},
				{Type: TokenTypeClosingBracket},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "array 2",
			raw:  "[true,false]",
			expected: []Token{
				{Type: TokenTypeOpeningBracket},
				{Type: TokenTypePrimitiveText, Value: "true"},
				{Type: TokenTypeComma},
				{Type: TokenTypePrimitiveText, Value: "false"},
				{Type: TokenTypeClosingBracket},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "object 1",
			raw:  `{"hello":true}`,
			expected: []Token{
				{Type: TokenTypeOpeningBrace},
				{Type: TokenTypeText, Value: "hello"},
				{Type: TokenTypeColon},
				{Type: TokenTypePrimitiveText, Value: "true"},
				{Type: TokenTypeClosingBrace},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 1",
			raw:  `"hi?"`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi?"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 2",
			raw:  `   "hi?"`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi?"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 3",
			raw:  `   "hi\n"`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi\\n"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 4 - escaping",
			raw:  `   "hi\""`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi\\\""},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 5 - escaping",
			raw:  `   "hi\n"`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi\\n"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "string 6",
			raw:  `   "hi(there]"`,
			expected: []Token{
				{Type: TokenTypeText, Value: "hi(there]"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "number 1",
			raw:  "123",
			expected: []Token{
				{Type: TokenTypePrimitiveNumber, Value: "123"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "number 2",
			raw:  "-123",
			expected: []Token{
				{Type: TokenTypePrimitiveNumber, Value: "-123"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "number 2",
			raw:  "   -01.23 ",
			expected: []Token{
				{Type: TokenTypePrimitiveNumber, Value: "-01.23"},
				{Type: TokenTypeEOF},
			},
		},
		{
			name: "nested 1",
			raw:  "{\"some\":true,\t\"arr\":[null\n,null]}",
			expected: []Token{
				{Type: TokenTypeOpeningBrace},
				{Type: TokenTypeText, Value: "some"},
				{Type: TokenTypeColon},
				{Type: TokenTypePrimitiveText, Value: "true"},
				{Type: TokenTypeComma},
				{Type: TokenTypeText, Value: "arr"},
				{Type: TokenTypeColon},
				{Type: TokenTypeOpeningBracket},
				{Type: TokenTypePrimitiveText, Value: "null"},
				{Type: TokenTypeComma},
				{Type: TokenTypePrimitiveText, Value: "null"},
				{Type: TokenTypeClosingBracket},
				{Type: TokenTypeClosingBrace},
				{Type: TokenTypeEOF},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := &lexer{r: bufio.NewReader(strings.NewReader(test.raw))}
			var tokens []Token
			for {
				token, err := lex.Token()
				if err != nil {
					assert.FailNow(t, "error lexing token", err)
				}
				assert.NoError(t, err)
				tokens = append(tokens, token)
				if token.Type == TokenTypeEOF {
					break
				}
			}
			assert.Equal(t, test.expected, tokens)
		})
	}
}

func TestLexer_LexPrimitiveText(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected string
	}{
		{
			name:     "simple file terminated string",
			raw:      "true",
			expected: "true",
		},
		{
			name:     "space terminated string 1",
			raw:      "true ",
			expected: "true",
		},
		{
			name:     "space terminated string 2",
			raw:      "true\t",
			expected: "true",
		},
		{
			name:     "space terminated string 3",
			raw:      "true\t",
			expected: "true",
		},
		{
			name:     "comma terminated string",
			raw:      "true,true,true",
			expected: "true",
		},
		{
			name:     "bracket terminated string",
			raw:      "true]",
			expected: "true",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := &lexer{r: bufio.NewReader(strings.NewReader(test.raw))}
			token, err := lex.lexPrimitiveText()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, token.Value)
		})
	}
}
