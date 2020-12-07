package main

import (
	"github.com/TuhinNair/durcov"
)

// Bot represents a message consuming and message producing conversational bot
type Bot struct {
	view *durcov.CovidBotView
}

func (b *Bot) respond(requestMessage string) string {
	return ""
}
