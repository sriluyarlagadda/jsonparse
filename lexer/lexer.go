package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

//type Lexer defines a lexical analyser
//for a JsonParser
type Lexer struct {
	reader *bufio.Reader
}

var escapeChars []rune = []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'}

//New initializes a new lexer which reads data from the given reader
func New(reader io.Reader) (*Lexer, error) {
	if reader == nil {
		return nil, errors.New("can't have empty reader")
	}
	lexer := &Lexer{}
	lexer.reader = bufio.NewReader(reader)
	return lexer, nil
}

// return the next valid token else return err
func (l *Lexer) NextToken() (Token, error) {

	isEnd := checkIfEndOfFile(l.reader)
	if isEnd {
		return Token{}, io.EOF
	}

	truncateWhiteSpace(l.reader)

	token, match, err := leftBracketToken(l.reader)
	if err != nil {
		return Token{}, err
	}

	if match {
		return token, nil
	}

	token, match, err = rightBracketToken(l.reader)
	if err != nil {
		return Token{}, err
	}

	if match {
		return token, nil
	}

	token, match, err = stringToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, nil
	}

	token, match, err = booleanToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, nil
	}

	token, match, err = nullToken(l.reader)
	fmt.Println("null token", token)
	if err != nil {
		return Token{}, err
	}

	if match {
		return token, nil
	}

	token, match, err = colonToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, nil
	}

	token, match, err = leftSquareBracketToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, nil
	}

	token, match, err = rightSquareBracketToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, nil
	}

	token, match, err = commaToken(l.reader)
	if err != nil {
		return Token{}, err
	}
	if match {
		return token, err
	}

	return Token{}, errors.New("unimplemented")
}

func checkIfEndOfFile(reader *bufio.Reader) bool {
	_, _, err := reader.ReadRune()

	if err == io.EOF {
		return true
	}

	rollBack(reader, 1)
	return false
}

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

func nullToken(reader *bufio.Reader) (Token, bool, error) {
	tokenString := make([]rune, 0)
	char, _, err := reader.ReadRune()
	if err != nil {
		printError(err, "nulltoken")
		os.Exit(1)
	}
	if char != 'n' {
		rollBack(reader, 1)
		return Token{}, false, nil
	}

	tokenString = append(tokenString, 'n')

	for i := 1; i <= 3; i++ {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				rollBack(reader, (i - 1))
				return Token{}, false, nil
			} else {
				printError(err, "nullToken")
				os.Exit(1)
			}
		}

		tokenString = append(tokenString, char)
	}
	fmt.Println("null token string", tokenString)

	fmt.Println("null token string", string(tokenString))
	if string(tokenString) == "null" {
		token := Token{Token_Null, "null"}
		return token, true, nil
	}
	rollBack(reader, 4)
	return Token{}, false, nil
}

func booleanToken(reader *bufio.Reader) (Token, bool, error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		printError(err, "boolean")
		os.Exit(1)
	}

	if char != 't' && char != 'f' {
		rollBack(reader, 1)
		return Token{}, false, nil
	}
	tokenString := make([]rune, 0)
	tokenString = append(tokenString, char)
	for i := 0; i < 4; i++ {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			printError(err, "booleanToken")
			os.Exit(1)
		}
		tokenString = append(tokenString, char)
	}

	state := 1
	for {
		switch state {
		case 0: //valid token
			return Token{Token_Bool, string(tokenString)}, true, nil
		case 1:
			if string(tokenString[0:4]) == "true" {
				rollBack(reader, 1)
				state = 0
				continue
			}

			if string(tokenString) == "false" {
				state = 0
				continue
			}

			rollBack(reader, 5)
			return Token{}, false, nil
		}
	}

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

func truncateWhiteSpace(reader *bufio.Reader) {
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return
			}
			printError(err, "truncateWhiteSpace")
			os.Exit(1)
		}

		if char != ' ' {
			reader.UnreadRune()
			break
		}
	}
}

func isUniCodeChar(char rune) bool {
	if char != '"' && char != '\\' && char != '/' {
		return true
	}
	return false
}
func isBackSlash(char rune) bool {
	if char == '\\' {
		return true
	}
	return false
}

func contains(slice []rune, search rune) bool {
	for _, currentRune := range slice {
		if currentRune == search {
			return true
		}
		return false
	}

	return false
}

func rollBack(reader *bufio.Reader, offset int) {
	for i := 0; i < offset; i++ {
		err := reader.UnreadRune()
		if err != nil {
			printError(err, "rollback")
		}
	}
}

func printError(err error, location string) {
	fmt.Println("error:", err)
	fmt.Println("Location:", location)
}
