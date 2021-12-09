package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

const PROMPT = ">> "

func ReplFromCommandLine(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	line := 0

	for {
		line++
		fmt.Printf("line %d %s", line, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		repl(scanner, out, env, line)
	}
}

func ReplFromFile(source *os.File, out io.Writer) {
	scanner := bufio.NewScanner(source)
	env := object.NewEnvironment()
	line := 0

	for {
		line++
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		repl(scanner, out, env, line)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "ERROR : "+msg+"\n")
	}
}

func repl(scanner *bufio.Scanner, out io.Writer, env *object.Environment, line int) {
	source := scanner.Text()
	l := lexer.New(source)
	p := parser.New(l, line)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
