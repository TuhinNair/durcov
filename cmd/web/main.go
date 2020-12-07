package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TuhinNair/durcov"

	"github.com/kevinburke/twilio-go"
)

type config struct {
	port            string
	twilioSID       string
	twilioAuthToken string
	dbURL           string
}

func loadConfig() *config {
	port := os.Getenv("PORT")
	port = ":" + port
	twilioSID := os.Getenv("TWILIO_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	dbURL := os.Getenv("DATABASE_URL")

	return &config{port, twilioSID, twilioAuthToken, dbURL}
}

func main() {
	config := loadConfig()
	pgxpool, err := durcov.GetPgxPool(config.dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxpool.Close()

	dataview := &durcov.CovidBotView{}
	dataview.SetDBConnection(pgxpool)
	bot := &Bot{dataview}

	twilioClient := twilio.NewClient(config.twilioSID, config.twilioAuthToken, nil)
	twilioBot := TwilioBot{twilioClient, bot}

	mux := http.NewServeMux()
	mux.HandleFunc("/whatsapp", twilioBot.handleWhatsapp)

	log.Printf("Starting server on port= %v", config.port)
	err = http.ListenAndServe(config.port, mux)
	if err != nil {
		log.Fatal(err)
	}

}
