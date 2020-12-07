package main

import (
	"github.com/TuhinNair/durcov"
)

type bot struct {
	view *durcov.CovidBotView
}

func (b *bot) respond(requestMessage string) (string, error) {
	return "", nil
}
