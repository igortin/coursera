package closers

import "fmt"

func main() {
	x := 997
	increment := func() int {
		x++
		return x
	}
	fmt.Println(increment())
	fmt.Println(increment())
}
