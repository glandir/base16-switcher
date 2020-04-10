package main

import (
	"github.com/alecthomas/kong"
)

var Cli struct {
	Update UpdateCmd `cmd:"" help:"Update templates and schemes."`
}

func main() {
	ctx := kong.Parse(&Cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
