package lexer

import "fmt"

type peekable struct {
	tokenizer Tokenizer
	peeked    *Token
	err       error
}

func (p *peekable) Peek() (Token, error) {
	if p.err != nil {
		return unknownToken, p.err
	}

	if p.peeked != nil {
		return *p.peeked, nil
	}

	token, err := p.Token()
	p.peeked = &token

	return token, err
}

func (p *peekable) Token() (Token, error) {
	if p.err != nil {
		return unknownToken, p.err
	}

	if p.peeked != nil {
		t := p.peeked
		p.peeked = nil
		return *t, nil
	}

	token, err := p.tokenizer.Token()
	if err != nil {
		err = fmt.Errorf("cannot lex token: %w", err)
		p.err = err
		return unknownToken, err
	}

	return token, nil
}
