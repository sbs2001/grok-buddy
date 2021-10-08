package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/logrusorgru/grokky"
)

func main() {
	exp := flag.String("expr", "def", "grok expression")
	lin := flag.String("log", "def", "log line")

	flag.Parse()
	linee, err := os.ReadFile(*lin)
	if err != nil {
		panic(err)
	}
	exprr, err := os.ReadFile(*exp)
	if err != nil {
		panic(err)
	}
	expr := string(exprr)
	line := string(linee)
	h := grokky.NewBase()
	n := len(expr)
	printInd := 0
	matchedPrinter := color.New(color.FgGreen).Add(color.Underline)
	unMatchedPrinter := color.New(color.FgRed).Add(color.Underline)
	foundMatch := false
	for i := n; i >= 0; i-- {
		tmpExpr := expr[:i]
		p, err := h.Compile(tmpExpr)
		// fmt.Println(p.String())
		if err != nil {
			continue
		}
		res := p.Parse(line)
		if len(res) != 0 {
			foundMatch = true
			fmt.Println("Grok Matched: ", tmpExpr)
			match := p.FindIndex([]byte(line))
			for ; printInd < len(line); printInd++ {
				if printInd == match[0] {
					for printInd < match[1] {
						matchedPrinter.Print(string(line[printInd]))
						printInd++
					}
				}
				if printInd == len(line) {
					break
				}
				unMatchedPrinter.Print(string(line[printInd]))
			}
			fmt.Println("\n")
			keys := make([]string, 0)
			for k := range res {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Println(fmt.Sprintf("evt.Parsed.%s = %s", k, res[k]))
			}
			break
		}
	}
	if !foundMatch {
		unMatchedPrinter.Println(line)
	}

	// fmt.Println(tmpExpr)
}

// }
