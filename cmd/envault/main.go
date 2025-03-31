package main

import (
	"fmt"
	"os"

	"github.com/uzulla/envault/internal/cli"
)

func main() {
	cli := cli.NewCLI()
	
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %s\n", err)
		os.Exit(1)
	}
}
