package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func parseConstantToken(reader *bufio.Reader, tokenType TokenType, tokenChar rune) (Token, bool, error) {
	char, _, err := reader.ReadRune()
	fmt.Println("parseConstant", string(char))
	if err != nil {
		if err == io.EOF {
			return Token{}, false, nil
		}
		printError(err, "parse Constant Token")
		os.Exit(1)
	}

	if char == tokenChar {
		token := Token{tokenType, string(tokenChar)}
		return token, true, nil
	}
	rollBack(reader, 1)
	return Token{}, false, nil
}

func leftSquareBracketToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_LBrace, '[')
}

func rightSquareBracketToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_RBrace, ']')
}

func colonToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_Colon, ':')
}

func commaToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_Comma, ',')
}

func rightBracketToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_RBracket, '}')
}

func leftBracketToken(reader *bufio.Reader) (Token, bool, error) {
	return parseConstantToken(reader, Token_LBracket, '{')
}
