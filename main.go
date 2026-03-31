package main

import (
	"fmt"
	"os"

	"github.com/yourusername/oze/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "oze: %v\n", err)
		os.Exit(1)
	}
}
