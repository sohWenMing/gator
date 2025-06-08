package utils

import (
	"fmt"
	"io"
)

func WriteLine(w io.Writer, input string) {
	fmt.Fprintln(w, input)
}
