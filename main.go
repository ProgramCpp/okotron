package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/okto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func main() {
	go auth()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	debugLevel, err := strconv.ParseBool(os.Getenv("ENABLE_DEBUG_LOGS"))
	if err != nil {
		log.Println("invalid value for config ENABLE_DEBUG_LOGS")
	}

	bot.Debug = debugLevel

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			reply := "hello " + "first_name"
			if update.Message.Text == "/auth" {
				reply = fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?response_type=code&client_id=%s&redirect_uri=http://localhost:3000/&scope=https://www.googleapis.com/auth/userinfo.profile&state=123&access_type=offline", os.Getenv("GOOGLE_CLIENT_ID"))
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func auth() {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		// todo: verify state
		// todo: use google auth package
		// https://pkg.go.dev/google.golang.org/api@v0.186.0/oauth2/v2

		code := r.URL.Query().Get("code")
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Println("error fetching google auth token " + err.Error())
		}

		idToken := token.Extra("id_token").(string) 

		// authenticate with okto
		oktoToken, err := okto.Authenticate(idToken)
		if err != nil {
			log.Println("authenticaiton with okto failed")
		}

		_ = oktoToken
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Println("auth server shutdown with error " + err.Error())
	}
}
