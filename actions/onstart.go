package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

var OnStart = SaveLastActionDecorate(func(app *application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		externalId := utils.ToString(message.Sender.ID)

		_, err := app.UserStore.Put(store.PutUser{
			ExternalId: externalId,
			Username:   message.Sender.Username,
			ChatId:     utils.ToString(message.Chat.ID),
		})

		if err != nil {
			return
		}

		log.Println(fmt.Sprintf("User %s (%s) created", message.Sender.Username, externalId))

		return
	}
})
