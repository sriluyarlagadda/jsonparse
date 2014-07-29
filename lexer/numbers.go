package lexer

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
)

func numberToken(reader *bufio.Reader) (Token, bool, error) {
	tokenString := make([]rune, 0)
	state := 1

Loop:
	for {
		switch state {

		case 0: //valid number
			token := Token{Token_Number, string(tokenString)}
			return token, true, nil
		case 1:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, false, nil
				}

				printError(err, "numberToken case 1")
				os.Exit(1)
			}

			if char == '-' {
				tokenString = append(tokenString, char)
				state = 7
				continue
			}

			if char == '0' {
				tokenString = append(tokenString, char)
				state = 6
				continue
			}

			if isNumeric(char) {
				tokenString = append(tokenString, char)
				state = 3
				continue
			}

			rollBack(reader, 1)
			break Loop

		case 2:

			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, false, errors.New("not a valid number")
				}
				printError(err, "number case 2")
				os.Exit(1)
			}

			if isNumeric(char) {
				tokenString = append(tokenString, char)
				state = 2
				continue
			}

			if char == 'e' || char == 'E' {
				tokenString = append(tokenString, char)
				state = 4
				continue
			}

		case 3:

			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					token := Token{Token_Number, string(tokenString)}
					return token, true, nil
				}
				break
			}
			if isNumeric(char) {
				tokenString = append(tokenString, char)
				state = 3
				continue
			}

			if char == '.' {
				tokenString = append(tokenString, char)
				state = 7
				continue
			}

			if char == 'e' || char == 'E' {
				tokenString = append(tokenString, char)
				state = 4
				continue
			}
		case 4:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, false, errors.New("exponent should have value")
				}
				printError(err, "numberToken case 4")
			}

			if char == '+' || char == '-' {
				tokenString = append(tokenString, char)
				state = 5
				continue
			}

			if isNumeric(char) {
				tokenString = append(tokenString, char)
				state = 5
				continue
			}

			return Token{}, false, errors.New("exponent should have value")
		case 5:
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{Token_Number, string(tokenString)}, true, nil
				}
				printError(err, "numberToken case 5")
			}

			if isNumeric(char) {
				tokenString = append(tokenString, char)
				state = 5
				continue
			}

			return Token{}, false, errors.New("should be atleast 1 digit")
		case 6:

			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{Token_Number, string(tokenString)}, true, nil
				}
				printError(err, "number token case 6")
				os.Exit(1)
			}

			if char == '.' {
				state = 2
				tokenString = append(tokenString, char)
				continue
			}

			if char == 'e' || char == 'E' {
				tokenString = append(tokenString, char)
				state = 4
				continue
			}

		case 7:

			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, false, errors.New("should be number after -")
				}

				printError(err, "numberToken case 7")
				os.Exit(1)
			}

			if char == '0' {
				state = 6
				tokenString = append(tokenString, char)
				continue
			}

			if isNumeric(char) {
				state = 3
				tokenString = append(tokenString, char)
				continue
			}

			return Token{}, false, errors.New("should be number after -")

		}
	}

	return Token{}, false, nil

}

func isNumeric(char rune) bool {
	stringNum := string(char)
	num, err := strconv.ParseInt(stringNum, 10, 8)
	if err != nil || (num < 1 && num > 9) {
		return false
	}

	return true
}
