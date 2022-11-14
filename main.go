package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andrewsjg/GoElite/internal/tui"
)

func main() {

	// TODO: Need a better way to encode and report the version.
	// TODO: Need better functionality to test the binary works properly
	const VER = "0.1.4"

	args := os.Args[1:]

	if len(args) >= 1 {
		arg0 := strings.ToUpper(args[0])

		if arg0 == "--VERSION" || arg0 == "--VER" {
			fmt.Println(VER)
			return
		}
	}

	tui.Start()

}
