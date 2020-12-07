package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kevinburke/twilio-go"
)

// TwilioBot represents a twilio client backed by a bot
type TwilioBot struct {
	client *twilio.Client
	bot    *Bot
}

type twilioRequest struct {
	to          string
	from        string
	requestBody string
}

type twilioResponse struct {
	to           string
	from         string
	responseBody string
}

func (tr *twilioRequest) toResponse(responseBody string) *twilioResponse {
	return &twilioResponse{to: tr.from, from: tr.to, responseBody: responseBody}
}

func (tr *twilioResponse) respond(twilioClient *twilio.Client) error {
	_, err := twilioClient.Messages.SendMessage(tr.from, tr.to, tr.responseBody, nil)
	return err
}

func (tb *TwilioBot) handleWhatsapp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}
	twilioRequestData, err := tb.parseRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = tb.respond(twilioRequestData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	} else {
		w.WriteHeader(200)
	}
	return
}

func (tb *TwilioBot) respond(reqData *twilioRequest) error {
	reqMsg := reqData.requestBody
	resMsg := tb.bot.respond(reqMsg)

	twilioResp := reqData.toResponse(resMsg)
	err := twilioResp.respond(tb.client)
	return err
}

func (tb *TwilioBot) parseRequest(r *http.Request) (*twilioRequest, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	to, ok := r.Form["To"]
	if !ok {
		log.Println("No To in form")
		return nil, err
	}
	if len(to) == 0 {
		log.Println("to is nil")
		return nil, errors.New("No to number")
	}
	toAddress := to[0]

	from, ok := r.Form["From"]
	if !ok {
		log.Println("No from in form")
		return nil, err
	}
	if len(from) == 0 {
		log.Println("No from number")
		return nil, errors.New("No from number")
	}
	fromAddress := from[0]

	body, ok := r.Form["Body"]
	if !ok {
		log.Println("No Body in form")
		return nil, err
	}
	if len(body) == 0 {
		log.Println("No message in body")
		return nil, errors.New("No message in body")
	}
	messageBody := body[0]

	twilioReq := twilioRequest{toAddress, fromAddress, messageBody}
	return &twilioReq, nil
}
