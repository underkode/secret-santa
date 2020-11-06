package actions

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

func OnText(app *application.ApplicationContext) func(message *tb.Message) {
	return func(message *tb.Message) {
		user := app.UserStore.FindByExternalId(utils.ToString(message.Sender.ID))

		if user == nil {
			fmt.Printf("User not found by %d", message.Sender.ID)
			return
		}

		lastAction := app.LastActionStore.FindByUserId(user.Id)

		switch lastAction.Action {
		case JOIN:
			Join(app, message, user)
		case PLAY_OUT:
			PlayOut(app, message, user)
		default:
			Default(app, message, user)
		}
	}
}

func Default(
	app *application.ApplicationContext,
	message *tb.Message,
	user *store.User,
) {
	_, err := app.Bot.Reply(message, message.Text)

	if err != nil {
		return
	}

	return
}
