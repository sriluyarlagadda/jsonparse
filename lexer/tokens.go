package lexer

import (
	"fmt"
)

type TokenType int

const (
	Token_LBracket TokenType = iota + 1
	Token_RBracket
	Token_Colon
	Token_Comma
	Token_LBrace
	Token_RBrace
	Token_String
	Token_Number
	Token_Bool
	Token_Null
	Token_Error
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token Type:%d Token Value: %s", t.Type, t.Value)
}

func getTokenTypeForTokens(char rune) TokenType {
	switch char {
	case '{':
		return Token_LBracket
	case '}':
		return Token_RBracket
	case ',':
		return Token_Comma
	case '[':
		return Token_LBrace
	case ']':
		return Token_RBrace
	}
	return 0
}
