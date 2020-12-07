package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/TuhinNair/durcov"
)

// Regex patterns to match valid commands.
var commandPattern = `(?P<command>^(?i)(CASES|DEATHS))`                                // Case insensitive match on either 'CASES' or 'DEATHS'
var countryCodePattern = `(?P<countryCode>(?i)((TOTAL)|[A-Z]{2})$)`                    // Case insensitive match on either 'TOTAL' or a two letter sequence from [A-Z]. Must be at the end of string.
var validBodyPattern = regexp.MustCompile(commandPattern + `\s+` + countryCodePattern) // Matches command and code delimited by 1+n whitespace

// Bot represents a message consuming and message producing conversational bot
type Bot struct {
	view *durcov.CovidBotView
}

type botError struct {
	err     error
	message string
	context []interface{}
}

func (b *Bot) handleBotError(botErr *botError) string {
	err := botErr.err
	responseMsg := botErr.message

	contextMsg := ""
	for i, info := range botErr.context {
		infoMsg := fmt.Sprintf("%d) %v\n", i, info)
		contextMsg = contextMsg + infoMsg
	}
	logMsg := fmt.Sprintf("\nError: %v\nContext: %s", err, contextMsg)
	log.Println(logMsg)

	err, ok := err.(*durcov.NoCountryMatchedError)
	if ok {
		responseMsg = "Sorry, that code doesn't match any countries I know."
	}
	return responseMsg
}

type requestType int

const (
	_Cases requestType = iota + 1
	_Deaths
)

type parsedRequest struct {
	Type requestType
	Code string
}

func (b *Bot) respond(requestMessage string) string {
	trimmedMsg, botErr := b.trimRequest(requestMessage)
	if botErr != nil {
		return b.handleBotError(botErr)
	}

	_, botErr = b.matchRequest(trimmedMsg)
	if botErr != nil {
		return b.handleBotError(botErr)
	}

	return ""

}

func (b *Bot) trimRequest(reqMsg string) (string, *botError) {
	trimmedReqMsg := strings.Trim(reqMsg, " ")
	if len(trimmedReqMsg) > 20 {
		requestLog := fmt.Sprintf("Message Too Long: %s", reqMsg)
		failedMessageCtxt := []interface{}{requestLog}
		botErr := &botError{
			errors.New("Request message too long"),
			"Sorry, that message is too long for me.",
			failedMessageCtxt,
		}
		return "", botErr
	}
	return trimmedReqMsg, nil
}

func (b *Bot) matchRequest(trimmedRequestMessage string) (*parsedRequest, *botError) {
	matches, botErr := b.subMatchRequest(trimmedRequestMessage)
	if botErr != nil {
		return nil, botErr
	}

	command, code := b.extractSubExp(matches)

	log.Printf("Got matched: Command: %s, Code: %s", command, code)
	switch command {
	case "CASES":
		return &parsedRequest{_Cases, code}, nil
	case "DEATHS":
		return &parsedRequest{_Deaths, code}, nil
	}

	unhandledCommandLog := fmt.Sprintf("Unhandled command: %s", command)
	unhandledCodeLog := fmt.Sprintf("Code in request: %v", code)
	failedMessageCtxt := []interface{}{
		unhandledCommandLog,
		unhandledCodeLog,
	}
	botErr = &botError{
		errors.New("Unhandled command"),
		"Oops, I've got myself confused :(",
		failedMessageCtxt,
	}
	return nil, botErr
}

func (b *Bot) subMatchRequest(reqMsg string) ([]string, *botError) {
	matches := validBodyPattern.FindStringSubmatch(reqMsg)
	if matches == nil {
		requestLog := fmt.Sprintf("Unmatchable Message: %s", reqMsg)
		failedMessageCtxt := []interface{}{requestLog}
		botErr := &botError{
			errors.New("Unmatched request message"),
			"Sorry, I'm not sure how to respond to that.",
			failedMessageCtxt,
		}
		return nil, botErr
	}
	return matches, nil
}

func (b *Bot) extractSubExp(matches []string) (command string, code string) {
	commandIdx := validBodyPattern.SubexpIndex("command")
	codeIdx := validBodyPattern.SubexpIndex("countryCode")

	command = matches[commandIdx]
	command = strings.ToUpper(command)
	code = matches[codeIdx]
	code = strings.ToUpper(code)
	return command, code
}
