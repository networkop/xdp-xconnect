package main

import (
	"fmt"
	"os"

	"github.com/networkop/xdp-xconnect/cmd"
)

var (
	GitCommit = "latest"
)

func main() {
	if err := cmd.Run(GitCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Error:\n%s\n", err.Error())
		os.Exit(1)
	}
}
