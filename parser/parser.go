package parser

import (
	"fmt"
	"strings"
)

// Parser represents a parser, including a scanner and the underlying raw input.
// It also contains a small buffer to allow for two unscans.
type Parser struct {
	s   *Lexer
	raw string
	buf TokenStack
}

// NewParser returns a new instance of Parser.
func NewParser(s string) *Parser {
	return &Parser{s: NewLexer(strings.NewReader(s)), raw: s}
}

// Parse takes the raw string and returns the root node of the AST.
func (p *Parser) Parse() (*Sql, error) {
	statement, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return statement, nil
}

func (p *Parser) parseExpression() (*Sql, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != SELECT {
		return nil, fmt.Errorf("we only scan SELECT statement. We found %s", lit)
	}

	// SELECT
	sqlSelect := make([]SqlSelect, 0)
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}
		if tok == FROM {
			break
		}
		if tok != STRING {
			return nil, fmt.Errorf("was exepecting a column name, or table.column")
		}

		s := SqlSelect{
			Column: lit,
		}
		tok, lit = p.scanIgnoreWhitespace()
		if tok == COMMA {
			sqlSelect = append(sqlSelect, s)
			continue
		}
		if tok == FROM || tok == EOF {
			sqlSelect = append(sqlSelect, s)
			break
		}
		if tok != DOT {
			return nil, fmt.Errorf("was exepecting a '.', i.e. table.column. We found %s (%d)", lit, tok)
		}
		tok, lit = p.scanIgnoreWhitespace()
		s.Table = s.Column
		s.Column = lit
		sqlSelect = append(sqlSelect, s)

		tok, lit = p.scanIgnoreWhitespace()
		if tok == FROM {
			sqlSelect = append(sqlSelect, s)
			break
		}
		if tok != COMMA {
			return nil, fmt.Errorf("was exepecting a ',' after a table.column")
		}
	}

	// FROM
	froms := make([]string, 0)
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}
		if tok == WHERE {
			break
		}
		if tok != STRING {
			return nil, fmt.Errorf("a table name was expected")
		}
		froms = append(froms, lit)

		tok, _ = p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}
		if tok == WHERE {
			break
		}
		if tok != COMMA {
			return nil, fmt.Errorf("a comma was expected")
		}
	}

	// WHERE
	where := ""
	for {
		tok, lit := p.scan()
		if tok == EOF {
			break
		}
		if tok == WS {
			where += " "
		} else {
			where += lit
		}
	}

	return &Sql{
		Select: sqlSelect,
		From:   froms,
		Where:  where,
	}, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.Len() != 0 {
		// Can ignore the error since it's not empty.
		tokenInf, _ := p.buf.Pop()
		return tokenInf.Token, tokenInf.Literal
	}

	// Otherwise read the next token from the scanner.
	tokenInf := p.s.Scan()
	tok, lit = tokenInf.Token, tokenInf.Literal
	return tok, lit
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

// unscan pushes the previously read tokens back onto the buffer.
func (p *Parser) unscan(tok TokenInfo) {
	p.buf.Push(tok)
}
