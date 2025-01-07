package main

import (
	"runtime"

	cmd "coingod/cmd/coingodcli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
