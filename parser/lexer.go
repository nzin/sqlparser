package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Lexer represents a lexical scanner.
type Lexer struct {
	r *bufio.Reader
}

// NewLexerFromString returns a Lexer for the provided string.
func NewLexerFromString(s string) *Lexer {
	return NewLexer(strings.NewReader(s))
}

// NewLexer returns a new instance of Lexer.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal Value.
func (s *Lexer) Scan() TokenInfo {
	// Read the next rune.
	ch := s.read()
	if ch == eof {
		return TokenInfo{EOF, ""}
	}

	switch {
	case ch == '.':
		return TokenInfo{DOT, string(ch)}
	case ch == '=':
		return TokenInfo{EQUAL, string(ch)}
	case ch == ',':
		return TokenInfo{COMMA, string(ch)}
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	default:
		s.unread()
		return s.scanKeyword()
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Lexer) scanWhitespace() TokenInfo {
	// Create a buffer and read the current character into it.
	var ch rune
	_ = s.read()
	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		ch = s.read()
		if ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}

	return TokenInfo{WS, ""}
}

// scanKeyword consumes the current rune and all contiguous text runes.
func (s *Lexer) scanKeyword() TokenInfo {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	quoted := false
	ch := s.read()
	if ch == '\'' {
		quoted = true
	} else {
		s.unread()
	}
	// Read every subsequent text character into the buffer.
	// Non-text characters and EOF will cause the loop to exit.
	for {
		ch := s.read()
		// Break if we hit EOF.
		if ch == eof {
			break
		}

		if ch == '\'' {
			break
		}
		// Break if we hit whitespace or a special char.
		if isWhitespace(ch) || isSpecialChar(ch) {
			s.unread()
			break
		}
		// Write the char into the buffer otherwise.
		buf.WriteRune(ch)
	}

	if quoted {
		return TokenInfo{QUOTED_STRING, buf.String()}
	}

	str := buf.String()
	if strings.EqualFold("SELECT", str) {
		return TokenInfo{SELECT, str}
	}
	if strings.EqualFold("FROM", str) {
		return TokenInfo{FROM, str}
	}
	if strings.EqualFold("WHERE", str) {
		return TokenInfo{WHERE, str}
	}

	return TokenInfo{STRING, buf.String()}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Lexer) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader, cannot unread twice sequentially.
func (s *Lexer) unread() {
	// Unread can error if we have previously not called read, this is not dangerous (no data mutation) and returning
	// error here would unnecessarily complicate the code.
	_ = s.r.UnreadRune()
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isSpecialChar(ch rune) bool {
	specialChar := []rune{'.', ',', '\'', '='}
	for _, char := range specialChar {
		if ch == char {
			return true
		}
	}
	return false
}

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
