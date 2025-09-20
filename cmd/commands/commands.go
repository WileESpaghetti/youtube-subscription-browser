package commands

import youtube_subscription_browser "github.com/WileESpaghetti/youtube-subscription-browser"

type Context struct {
	Verbose   bool   `help:"Enable verbose mode."`
	Database  string `help:"Database file." default:"./youtube.sqlite"`
	TokenFile string `help:"OAuth authorization token." default:""`
}

func NewContext() *Context {
	return &Context{
		Verbose:   false,
		Database:  "./youtube.sqlite",
		TokenFile: youtube_subscription_browser.DefaultTokenFile,
	}
}
