package lexer

type Lexer struct {
	input        string
	position     int  //入力における現在の位置(現在の文字を指し示す)
	readPosition int  //これから読み込む位置(現在の文字の次)
	ch           byte //現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() //l.ch,position,readPositionを初期化
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 //ASCII:"NUL"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
