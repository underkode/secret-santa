package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

type Action func(app *application.ApplicationContext) func(message *tb.Message)

func SaveLastActionDecorate(next Action) Action {
	return func(app *application.ApplicationContext) func(message *tb.Message) {
		return func(message *tb.Message) {
			user, _ := app.UserStore.Put(store.PutUser{
				ExternalId: utils.ToString(message.Sender.ID),
				Username:   message.Sender.Username,
				ChatId:     utils.ToString(message.Chat.ID),
			})

			if user == nil {
				fmt.Printf("User not found by %d", message.Sender.ID)
				return
			}

			_, err := app.LastActionStore.Put(store.PutLastAction{
				Action:    message.Text,
				OwnerUser: user,
			})

			if err != nil {
				return
			}

			next(app)(message)
		}
	}
}
