package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	TokenTypeUnknown TokenType = iota
	TokenTypeOpeningBrace
	TokenTypeClosingBrace
	TokenTypeOpeningBracket
	TokenTypeClosingBracket
	TokenTypeComma
	TokenTypeColon

	// TokenTypeText is a token of escaped text in double quotes
	TokenTypeText

	// TokenTypePrimitiveText is a token of unescaped, lower-case alpha text
	TokenTypePrimitiveText

	// TokenTypePrimitiveNumber is a token of a json number
	TokenTypePrimitiveNumber

	TokenTypeEOF
)

const delimiters = ",{}[]:"

type Token struct {
	Type  TokenType
	Value string
}

var (
	unknownToken        = Token{Type: TokenTypeUnknown}
	eofToken            = Token{Type: TokenTypeEOF}
	openingBraceToken   = Token{Type: TokenTypeOpeningBrace}
	closingBraceToken   = Token{Type: TokenTypeClosingBrace}
	openingBracketToken = Token{Type: TokenTypeOpeningBracket}
	closingBracketToken = Token{Type: TokenTypeClosingBracket}
	colonToken          = Token{Type: TokenTypeColon}
	commaToken          = Token{Type: TokenTypeComma}
)

type Tokenizer interface {
	// Token reads the next token, if there is any error it will keep returning the same error
	// on subsequent calls
	Token() (Token, error)
}

type Lexer interface {
	Tokenizer

	Peek() (Token, error)
}

type lexer struct {
	r *bufio.Reader
}

func NewLexer(r io.Reader) Lexer {
	return &peekable{
		tokenizer: &lexer{
			r: bufio.NewReader(r),
		},
	}
}

func (l *lexer) Token() (Token, error) {
	err := l.skipWhitespace()
	if err != nil {
		if err == io.EOF {
			return eofToken, nil
		}
		return unknownToken, err
	}

	r, _, err := l.r.ReadRune()
	switch r {
	case '"':
		return l.lexText()
	case '{':
		return openingBraceToken, nil
	case '}':
		return closingBraceToken, nil
	case '[':
		return openingBracketToken, nil
	case ']':
		return closingBracketToken, nil
	case ':':
		return colonToken, nil
	case ',':
		return commaToken, nil
	}

	err = l.r.UnreadRune()
	if err != nil {
		return unknownToken, err
	}

	if isNumberLiteral(r) {
		return l.lexPrimitiveNumber()
	}

	if isAlpha(r) {
		return l.lexPrimitiveText()
	}

	return unknownToken, fmt.Errorf("unrecognized token: '%c'", r)
}

func (l *lexer) skipWhitespace() error {
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			return l.r.UnreadRune()
		}
	}
}

func (l *lexer) lexPrimitiveNumber() (Token, error) {
	v, err := l.readPrimitiveNumber()
	if err != nil {
		return unknownToken, err
	}

	// TODO: Validate literal?

	return Token{Type: TokenTypePrimitiveNumber, Value: v}, nil
}

func (l *lexer) readPrimitiveNumber() (string, error) {
	var s strings.Builder
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return s.String(), nil
			}
			return "", err
		}

		if isDelimiter(r) {
			return s.String(), l.r.UnreadRune()
		}

		_, err = s.WriteRune(r)
		if err != nil {
			return "", err
		}
	}
}

func (l *lexer) lexPrimitiveText() (Token, error) {
	v, err := l.readPrimitiveText()
	if err != nil {
		return unknownToken, err
	}

	return Token{Type: TokenTypePrimitiveText, Value: v}, nil
}

func (l *lexer) readPrimitiveText() (string, error) {
	var s strings.Builder
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return s.String(), nil
			}
			return "", err
		}

		if isDelimiter(r) {
			return s.String(), l.r.UnreadRune()
		}

		if !isAlpha(r) {
			return "", fmt.Errorf("invalid primitive text: '%s%c'", s.String(), r)
		}

		_, err = s.WriteRune(r)
		if err != nil {
			return "", err
		}
	}
}

func (l *lexer) lexText() (Token, error) {
	v, err := l.readText()
	if err != nil {
		return unknownToken, err
	}

	return Token{Type: TokenTypeText, Value: v}, nil
}

func (l *lexer) readText() (string, error) {
	var s strings.Builder
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("unexpected EOF in text")
			}
			return "", err
		}

		if r == '"' {
			return s.String(), nil
		}
		if r == '\\' {
			_, err = s.WriteRune(r)
			if err != nil {
				return "", err
			}

			r, _, err = l.r.ReadRune()
			if err != nil {
				return "", err
			}
		}

		_, err = s.WriteRune(r)
		if err != nil {
			return "", err
		}
	}
}

func isDelimiter(r rune) bool {
	return unicode.IsSpace(r) || strings.ContainsRune(delimiters, r)
}

func isNumberLiteral(r rune) bool {
	return unicode.IsDigit(r) || r == '-'
}

func isAlpha(r rune) bool {
	return 'a' <= r && 'z' >= r
}
