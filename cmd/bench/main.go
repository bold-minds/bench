package main

import (
	"os"

	"github.com/bold-minds/bench/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
