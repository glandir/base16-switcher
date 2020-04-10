package main

import (
	"github.com/alecthomas/kong"
)

var Cli struct {
	Update UpdateCmd `cmd:"" help:"Update templates and schemes."`
	List ListCmd `cmd:"" help:"List available color schemes."`
	Apply ApplyCmd `cmd:"" help:"Apply the named color scheme or use the default if none is specified."`
}

func main() {
	ctx := kong.Parse(&Cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
