package parser

// Token represents a lexical token.
type Token int

// TokenInfo stores relevant information about the token during scanning.
type TokenInfo struct {
	Token   Token
	Literal string
}

// TokenLookup is a map, useful for printing readable names of the tokens.
var TokenLookup = map[Token]string{
	OTHER:         "OTHER",
	EOF:           "EOF",
	WS:            "WS",
	STRING:        "STRING",
	QUOTED_STRING: "QUOTED_STRING",
	DOT:           ".",
	COMMA:         ",",
	EQUAL:         "=",
	SELECT:        "SELECT",
	FROM:          "FROM",
	WHERE:         "WHERE",
}

// String prints a human readable string name for a given token.
func (t Token) String() (print string) {
	return TokenLookup[t]
}

// Declare the tokens here.
const (
	// Special tokens
	// Iota simply starts and integer count
	OTHER Token = iota
	EOF
	WS

	// Main literals
	STRING
	QUOTED_STRING

	DOT
	EQUAL
	COMMA

	// Operators
	PLUS
	MINUS
	MULTIPLY
	DIVIDE

	// SQL specific
	SELECT
	FROM
	WHERE
)
