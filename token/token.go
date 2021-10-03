package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//識別子+リテラル
	IDENT = "IDENT" //add,foobar,s,y,...
	INT   = "INT"   //1343456

	//デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	キーワード
	FUNCTIOIN = "FUNCTION"
	LET       = "LET"
)
