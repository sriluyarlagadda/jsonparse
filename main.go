package main

import (
	"strings"

	"github.com/sriluyarlagadda/jsonparse/lexer"
)

func main() {
	var inputString string = "hai"
	reader := strings.NewReader(inputString)
	jsonLexer := lexer.New(reader)

	token, err := jsonLexer.NextToken()

}
