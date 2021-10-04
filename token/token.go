package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIndent(indent string) TokenType {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return IDENT
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//識別子+リテラル
	IDENT = "IDENT" //add,foobar,s,y,...
	INT   = "INT"   //1343456

	//演算子
	ASSIGN = "="
	PLUS   = "+"

	//デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
