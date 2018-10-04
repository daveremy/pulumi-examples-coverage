package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func main() {

	outputFlag := flag.String("output", "console", "Output options: 'console' or 'file'.")

	examples := getExamples()
	examplesWithTests := getExamplesWithTests()

	if *outputFlag == "console" {
		coverageReportToConsole(examples, examplesWithTests)
	} else {
		log.Fatal("output option 'file' not implemented yet")
	}
}

func coverageReportToConsole(examples []string, withTests []string) {

	longestExampleName := longestExampleName(examples)
	color.Cyan("%sTest Exists", addPad("Example Name", longestExampleName))
	color.Cyan("%s-----------", addPad("------------", longestExampleName))
	for _, e := range examples {
		if contains(e, withTests) {
			color.Cyan("%sYes", addPad(e, longestExampleName))
		} else {
			color.Red("%sNO", addPad(e, longestExampleName))
		}
	}
	color.Unset()
}

func addPad(s string, margin int) string {
	padLength := margin + 4 - len(s)
	return s + strings.Repeat(" ", padLength)
}

func longestExampleName(examples []string) int {
	result := 0
	for _, e := range examples {
		if len(e) > result {
			result = len(e)
		}
	}
	return result
}

// Parse the examples test file source code in order to get the implemented
//  tests.
func getExamplesWithTests() []string {
	var result []string
	goRoot := os.Getenv("GOPATH")
	examplesTestSrc := filepath.Join(goRoot, "/src/github.com/pulumi/examples/misc/test/", "examples_test.go")

	fset := token.NewFileSet()
	// Parse the source file here ...
	f, err := parser.ParseFile(fset, examplesTestSrc, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// parse the examples test
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			// TODO: Change to use integration.ProgramOptions rather than this Join hack.
			// Relies on the fact that the examples test join the test path together
			// and the last string is the example name
			callHasJoin := false
			switch v := x.Fun.(type) {
			case *ast.SelectorExpr:
				if v.Sel.Name == "Join" {
					callHasJoin = true
				}
			}
			if callHasJoin {
				switch s := x.Args[3].(type) {
				case *ast.BasicLit:
					result = append(result, s.Value[1:len(s.Value)-1])
				}
			}
		}
		return true
	})
	return result
}

func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getExamples() []string {
	var result []string
	goRoot := os.Getenv("GOPATH")
	examplesDir := filepath.Join(goRoot, "/src/github.com/pulumi/examples/")
	files, err := ioutil.ReadDir(examplesDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			excludedDirs := []string{".git", "misc", "~."}
			if !contains(f.Name(), excludedDirs) {
				result = append(result, f.Name())
			}
		}
	}
	return result
}
