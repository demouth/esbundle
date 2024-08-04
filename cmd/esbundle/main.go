package main

import (
	"os"

	"github.com/demouth/esbundle/pkg/cli"
)

func main() {
	osArgs := os.Args[1:]
	cli.Run(osArgs)
}
