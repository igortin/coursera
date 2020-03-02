package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var (
	testOk = `1
2
2
3
4
5
5
5
5
5`

	testOkResult = `1
2
3
4
5
`
)

func TestOK(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)
	err := uniq(in, out)
	// fmt.Println(out.String())

	// unit test internal
	if err != nil {
		t.Errorf("Error %v", err.Error())
	}

	// unit test match string good
	if out.String() != testOkResult {
		t.Errorf("Not matched %v\n%v", out.String(), testOkResult)
	}
}
