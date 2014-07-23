package lexer

import (
	"fmt"
	"strings"
	"testing"
)

var testMap map[string]Token = make(map[string]Token)
var multiTokenMap map[string][]Token = make(map[string][]Token)

func init() {
	testMap["\"hai\""] = Token{Token_String, "\"hai\""}
	testMap["\"\""] = Token{Token_String, "\"\""}
	testMap["\"\rasdf\""] = Token{Token_String, "\"\rasdf\""}
	testMap[" \" asd\" "] = Token{Token_String, "\" asd\""}

	multiTokenMap["false \"jfk\n\" "] = []Token{
		Token{Token_Bool, "false"},
		Token{Token_String, "\"jfk\n\""}}

	multiTokenMap["false { \"sd f\"} \" \" "] = []Token{
		Token{Token_Bool, "false"},
		Token{Token_LBracket, "{"},
		Token{Token_String, "\"sd f\""},
		Token{Token_RBracket, "}"},
		Token{Token_String, "\" \""},
	}

	multiTokenMap["false : { } \"null\" null [, true"] = []Token{
		Token{Token_Bool, "false"},
		Token{Token_Colon, ":"},
		Token{Token_LBracket, "{"},
		Token{Token_RBracket, "}"},
		Token{Token_String, "\"null\""},
		Token{Token_Null, "null"},
		Token{Token_LBrace, "["},
		Token{Token_Comma, ","},
		Token{Token_Bool, "true"},
	}
}

func TestStringToken(t *testing.T) {

	for key, testToken := range testMap {
		var inputString string = key
		reader := strings.NewReader(inputString)
		jsonLexer, _ := New(reader)
		token, _ := jsonLexer.NextToken()

		if token.Type != testToken.Type {
			t.Error("token type not string", token.Value)
		}

		if token.Value != testToken.Value {
			t.Error("value should be ", testToken.Value, "but was", token.Value)
		}

	}
}

func TestForwardSlashWithoutEscape(t *testing.T) {
	var inputString string = "\"/\""
	reader := strings.NewReader(inputString)
	jsonLexer, _ := New(reader)
	_, err := jsonLexer.NextToken()

	if err == nil {
		t.Error("should produce error")
	}
}

func TestBracketAndSlash(t *testing.T) {
	inputString := "{ \"asdf\""
	reader := strings.NewReader(inputString)
	jsonLexer, _ := New(reader)
	token, _ := jsonLexer.NextToken()

	leftBracket := Token{Token_LBracket, string('{')}
	if token != leftBracket {
		t.Error("expected ", string('{'), " but got", token.Value)
	}

	if token.Type != Token_LBracket {
		t.Error("expected ", Token_LBracket, " but got", token.Type)
	}
}

func TestMultiToken(t *testing.T) {
	for key, value := range multiTokenMap {
		inputString := key
		reader := strings.NewReader(inputString)
		jsonLexer, _ := New(reader)

		for _, tokenVal := range value {
			token, err := jsonLexer.NextToken()
			fmt.Println(err)
			fmt.Println("tokenVal:", tokenVal, "got token:", token)

			if token.Type != tokenVal.Type {
				t.Error("expected type ", string(tokenVal.Type), " but got type", string(token.Type))

			}

			if token.Value != tokenVal.Value {
				t.Error("expected val", string(tokenVal.Value), " but got val", string(token.Value))
			}
		}
	}
}

func TestMultiTokenWithWrongBool(t *testing.T) {
	inputString := "false nil \"jfk\n\" "
	reader := strings.NewReader(inputString)
	jsonLexer, _ := New(reader)
	token, err := jsonLexer.NextToken()

	if token.Value != "false" {
		t.Error("expected token", "false", " got", token.Value)
	}

	_, err = jsonLexer.NextToken()
	if err == nil {
		t.Error("expected err ")
	}
}
