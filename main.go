package main

import (
	"fmt"
	"os"

	"github.com/nzin/sqlparser/parser"
)

func main() {
	expr := ""
	for _, a := range os.Args[1:] {
		expr += a + " "
	}

	p := parser.NewParser(expr)
	statement, err := p.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Printf("SELECT ")
	firstTable := true
	for _, s := range statement.Select {
		if firstTable == false {
			fmt.Printf(",")
		} else {
			firstTable = false
		}
		if s.Table == "" {
			fmt.Printf(s.Column)
		} else {
			fmt.Printf("%s.%s", s.Table, s.Column)
		}
	}

	if len(statement.From) > 0 {
		fmt.Printf(" FROM ")
		firstFrom := true
		for _, f := range statement.From {
			if firstFrom == false {
				fmt.Printf(",")
			} else {
				firstFrom = false
			}
			fmt.Printf(f)
		}
		if len(statement.Where) > 0 {
			fmt.Printf(" WHERE%s", statement.Where)
		}
	}
	fmt.Println()
}
