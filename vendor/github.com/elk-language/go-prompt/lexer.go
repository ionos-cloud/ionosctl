package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// Lexer is a streaming lexer that takes in a piece of text
// and streams tokens with the Next() method
type Lexer interface {
	Init(string) // Reset the lexer's state and initialise it with the given input.
	// Next returns the next Token and a bool flag
	// which is false when the end of input has been reached.
	Next() (Token, bool)
}

// Token is a single unit of text returned by a Lexer.
type Token interface {
	Color() Color // Color of the token's text
	BackgroundColor() Color
	DisplayAttributes() []DisplayAttribute
	FirstByteIndex() istrings.ByteNumber // Index of the last byte of this token
	LastByteIndex() istrings.ByteNumber  // Index of the last byte of this token
}

// SimpleToken as the default implementation of Token.
type SimpleToken struct {
	color             Color
	backgroundColor   Color
	displayAttributes []DisplayAttribute
	lastByteIndex     istrings.ByteNumber
	firstByteIndex    istrings.ByteNumber
}

type SimpleTokenOption func(*SimpleToken)

func SimpleTokenWithColor(c Color) SimpleTokenOption {
	return func(t *SimpleToken) {
		t.color = c
	}
}

func SimpleTokenWithBackgroundColor(c Color) SimpleTokenOption {
	return func(t *SimpleToken) {
		t.backgroundColor = c
	}
}

func SimpleTokenWithDisplayAttributes(attrs ...DisplayAttribute) SimpleTokenOption {
	return func(t *SimpleToken) {
		t.displayAttributes = attrs
	}
}

// Create a new SimpleToken.
func NewSimpleToken(firstIndex, lastIndex istrings.ByteNumber, opts ...SimpleTokenOption) *SimpleToken {
	t := &SimpleToken{
		firstByteIndex: firstIndex,
		lastByteIndex:  lastIndex,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Retrieve the text color of this token.
func (t *SimpleToken) Color() Color {
	return t.color
}

// Retrieve the background color of this token.
func (t *SimpleToken) BackgroundColor() Color {
	return t.backgroundColor
}

// Retrieve the display attributes of this token eg. bold, underline.
func (t *SimpleToken) DisplayAttributes() []DisplayAttribute {
	return t.displayAttributes
}

// The index of the last byte of the lexeme.
func (t *SimpleToken) LastByteIndex() istrings.ByteNumber {
	return t.lastByteIndex
}

// The index of the first byte of the lexeme.
func (t *SimpleToken) FirstByteIndex() istrings.ByteNumber {
	return t.firstByteIndex
}

// LexerFunc is a function implementing
// a simple lexer that receives a string
// and returns a complete slice of Tokens.
type LexerFunc func(string) []Token

// EagerLexer is a wrapper around LexerFunc that
// transforms an eager lexer which produces an
// array with all tokens at once into a streaming
// lexer compatible with go-prompt.
type EagerLexer struct {
	lexFunc      LexerFunc
	tokens       []Token
	currentIndex int
}

// Create a new EagerLexer.
func NewEagerLexer(fn LexerFunc) *EagerLexer {
	return &EagerLexer{
		lexFunc: fn,
	}
}

// Initialise the lexer with the given input.
func (l *EagerLexer) Init(input string) {
	l.tokens = l.lexFunc(input)
	l.currentIndex = 0
}

// Return the next token and true if the operation
// was successful.
func (l *EagerLexer) Next() (Token, bool) {
	if l.currentIndex >= len(l.tokens) {
		return nil, false
	}

	result := l.tokens[l.currentIndex]
	l.currentIndex++
	return result, true
}
