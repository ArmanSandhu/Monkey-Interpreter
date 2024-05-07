package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/armansandhu/monkey_interpreter/lexer"
	"github.com/armansandhu/monkey_interpreter/token"
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

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
