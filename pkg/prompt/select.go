package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Select blocks until a valid integer is selected
func Select(options ...string) int {
	fmt.Println("\nmake a selection:")
	for i, option := range options {
		fmt.Printf("\t%d: %s\n", i+1, option)
	}

	reader := bufio.NewReader(os.Stdin)
	b, _, err := reader.ReadLine()

	if err != nil {
		fmt.Println(err.Error())
		Select(options...)
	}

	i, err := strconv.Atoi(string(b))
	if err != nil {
		fmt.Println(err.Error())
		Select(options...)
	}

	if i <= 0 || i > len(options) {
		fmt.Println(ErrInvalidSelection.Error())
		Select(options...)
	}

	return i
}
