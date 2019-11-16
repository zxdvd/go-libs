package debug

import (
	"testing"
	"fmt"
)

func TestGetSymbolOffset(t *testing.T) {
	offset, err := GetSymbolOffset("/bin/bash", "readline")
	fmt.Println(offset, err)
	if err != nil {
		panic(err)
	}
}
