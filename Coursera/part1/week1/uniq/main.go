package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	err := uniq(os.Stdin, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func uniq(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string

	for in.Scan() {
		line := in.Text()
		if line < prev {
			return fmt.Errorf("file not sorted")
		}

		if line == prev {
			continue
		}
		prev = line
		fmt.Fprintln(output, line)
	}
	return nil
}
