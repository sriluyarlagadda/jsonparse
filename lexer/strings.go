package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func stringToken(reader *bufio.Reader) (Token, bool, error) {
	count := 0
	tokenString := make([]rune, 0)
	state := 1
L:
	for {
		switch state {
		case 0: //valid string
			return Token{Token_String, string(tokenString)}, true, nil
		case 1:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return Token{}, false, err
			}
			count++

			if char == '"' {
				tokenString = append(tokenString, char)
				state = 2
				continue
			}
			rollBack(reader, count)
			break L //exist the entire for loop
		case 2:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					break L
				}
				printError(err, "Striing token case 2")
			}

			if char == '"' {
				tokenString = append(tokenString, char)
				count++
				state = 0
				continue
			}

			if isBackSlash(char) {
				tokenString = append(tokenString, char)
				count++
				state = 3
				continue
			}

			if isUniCodeChar(char) {
				tokenString = append(tokenString, char)
				count++
				state = 2
				continue
			}

			errString := fmt.Sprintf("illegal character %s", string(char))
			err = errors.New(errString)
			printError(err, "String Token: case 2")
			return Token{}, false, err

		case 3:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				printError(err, "String Token case 3")
				return Token{}, false, err
			}
			if !contains(escapeChars, char) {
				err := errors.New("cannot have any other character after \\")
				printError(err, "stringToken")
			}

			tokenString = append(tokenString, char)
			count++
			state = 2
		}

	}
	return Token{}, false, nil
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
