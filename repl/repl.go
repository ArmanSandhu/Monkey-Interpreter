package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/armansandhu/monkey_interpreter/lexer"
	"github.com/armansandhu/monkey_interpreter/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scan := scanner.Scan()
		if !scan {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		parse := parser.New(lex)

		program := parse.ParseProgram()
		if len(parse.Errors()) != 0 {
			printParseErrors(out, parse.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
