package telegram

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/programcpp/oktron/db"
	"github.com/programcpp/oktron/google"
	"github.com/programcpp/oktron/okto"
)

func Auth(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	// dbKey := fmt.Sprintf("%d", id)
	// token := db.Get(dbKey)
	// // TODO: do proper error check. check error from db call
	// if token != "" {
	// 	Send(bot, update, "Oktron already authorized.")
	// 	return
	// }

	deviceCode, err := google.GetDeviceCode()
	if err != nil {
		log.Printf("error encountered when saving token id %d", id)
		Send(bot, update, "error authorizing oktron. try again.")
		return
	}
	reply := ""
	reply += fmt.Sprintf("visit [google authorization page](%s) to authorize oktron and enter device code %s. \n", deviceCode.VerificationUrl, deviceCode.UserCode)
	reply += "return to oktron chat when done"
	Send(bot, update, reply)

	for i := 0; i*deviceCode.Interval <= deviceCode.ExpiresIn; i++ {
		time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
		googleToken, err := google.PollAuthorization(deviceCode.DeviceCode)
		if errors.Is(err, google.ErrAuthorizationPending) {
			continue
		} else if err != nil {
			log.Printf("error encountered when polling for token id for user %d", id)
			Send(bot, update, "error authorizing oktron. try again.")
			break
		} else {

			googleIDTokenKey := fmt.Sprintf("google_id_token_%d", id) // TODO: expire tokens
			err = db.Save(googleIDTokenKey, googleToken.IdToken)
			if err != nil {
				log.Printf("error encountered when saving google token id %d. %s", id, err.Error())
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			token, err := okto.Authenticate(googleToken.IdToken)
			if err != nil {
				log.Println("error authentication to Okto. " + err.Error())
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			tokenKey := fmt.Sprintf("okto_token_%d", id) // TODO: expire tokens
			err = db.Save(tokenKey, token)
			if err != nil {
				log.Printf("error encountered when saving token id %d. %s", id, err.Error())
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			resp, err := SendWithForceReply(bot, update, "almost done! set your - 6 digit - PIN to finish setup", true)
			if err != nil {
				log.Printf("error encountered when sending bot message. %s", err.Error())
				// Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			// commands in progress
			// multi-step commands are sequenced with message id's
			// save the next command in the sequence
			// the flow works only if the user replies to the message sent by the bot
			// this allows the bot to determine the next command based on the user flow instead of the user manually selecting the commands. this improves the UX and simplifies bot usage
			// TODO: document this user flow in a ADR
			messageKey := fmt.Sprintf("message_%d", resp.MessageID)
			err = db.Save(messageKey, PIN)
			if err != nil {
				log.Printf("error encountered when saving token id %d. %s", id, err.Error())
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}
			break
		}
	}
}
