package main

import (
	"github.com/WileESpaghetti/youtube-subscription-browser/cmd/commands"
	"github.com/alecthomas/kong"
)

var cli struct {
	commands.Context

	Auth   commands.AuthCmd   `cmd:"" help:"Authenticate with YouTube Data API"`
	Import commands.ImportCmd `cmd:"" help:"Import"`
	InitDB commands.InitDBCmd `cmd:"" help:"init-db"`
}

func main() {
	ctx := kong.Parse(&cli, kong.ShortUsageOnError())
	err := ctx.Run(&cli.Context)
	ctx.FatalIfErrorf(err)
}
