package telegram

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
)

var (
	// TODO: use okto /aupported_tokens and /supported_networks api's
	// do not hardcode networks and tokens
	// for now, all networks returnd by /supported_networks do not work. ex: solana, osmosis
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. oktron is network agnostic
	SUPPORTED_TOKENS = []string{"APT", "ETH", "MATIC", "USDC", "USDT"}
	SUPPORTED_NETWORKS   = map[string][]string{
		"APT":   {"APTOS"},
		"ETH":   {"BASE"},
		"MATIC": {"POLYGON"},
		"USDC":  {"BASE", "POLYGON"},
		"USDT":  {"POLYGON", "APTOS"},
	}
)

type SwapRequest struct {
	FromNetwork string `json:"from_network"`
	ToNetwork   string `json:"to_network"`
	FromToken   string `json:"from_token"`
	ToToken     string `json:"to_token"`
}

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: separate UI componets from backend
	keyboardButtons := []tgbotapi.InlineKeyboardButton{}

	// TODO: see user portfolio for a customized user experience
	for _, n := range SUPPORTED_TOKENS {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var tokenKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)
	// TODO: handle error
	resp, _ := bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      update.Message.Chat.ID,
			ReplyMarkup: tokenKeyboard,
		},
		Text: "select the source token",
	})

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err := db.Save(messageKey, CMD_SWAP_CMD_SOURCE_TOKEN)
	if err != nil {
		log.Printf("error encountered when saving swap message id. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapSourceToken(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	fromToken := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID
	requestKey := fmt.Sprintf("swap_%d", id)
	var buf bytes.Buffer
	// TODO:handle error
	_ = json.NewEncoder(bufio.NewWriter(&buf)).Encode(SwapRequest{FromToken: fromToken})
	err := db.Save(requestKey, buf.String())
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_NETWORKS[fromToken] {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "", networkKeyboard)
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.Save(messageKey, CMD_SWAP_CMD_SOURCE_NETWORk)
	if err != nil {
		log.Printf("error encountered when saving swap network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}
