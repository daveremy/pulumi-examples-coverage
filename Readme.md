## Pulumi Examples Test Coverage

Simple go program that lists the examples in the main pulumi example repo (https://github.com/pulumi/examples) and whether or not the examples has a test.

The program requires that the example repo is in `$GOPATH/src/github.com/pulumi/examples`.

The program uses a couple of heuristics that may not ultimately always hold true.

1. To get the list of examples it reads the top level directories in the examples repo and assumes that each example has it's own directory (it filters out directories like `.git` and `misc`).

1. To get the list of tests it reads the `misc/test/examples_test.go` source code into an abstract syntax tree and grabs the tests.  (* Note, if a test is commented out it will not find it).

