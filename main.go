package main

import (
	"github.com/spf13/pflag"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
	"underkode.ru/secret-santa/actions"
	"underkode.ru/secret-santa/application"
	"underkode.ru/secret-santa/store"
	"underkode.ru/secret-santa/utils"
)

func newUserStore(workDirectory string) *store.UserStore {
	return utils.CheckAndReturn(store.NewUserStore(workDirectory + ".users")).(*store.UserStore)
}

func newRoundStore(workDirectory string) *store.RoundStore {
	return utils.CheckAndReturn(store.NewRoundStore(workDirectory + ".rounds")).(*store.RoundStore)
}

func newLastActionStore(workDirectory string) *store.LastActionStore {
	return utils.CheckAndReturn(store.NewLastActionStore(workDirectory + ".lastactions")).(*store.LastActionStore)
}

func newRoundParticipantStore(workDirectory string) *store.RoundParticipantStore {
	return utils.CheckAndReturn(store.NewRoundParticipantStore(workDirectory + ".roundparticipants")).(*store.RoundParticipantStore)
}

func newBot(tgToken string) *tb.Bot {
	return utils.CheckAndReturn(tb.NewBot(tb.Settings{
		Verbose: true,
		Token:   tgToken,
		Poller:  tb.NewMiddlewarePoller(&tb.LongPoller{Timeout: 10 * time.Second}, utils.HasMessageFromPrivateChat),
	})).(*tb.Bot)
}

func main() {
	workDirectory := pflag.StringP("work-directory", "d", "", "work directory")
	tgToken := pflag.StringP("tg-token", "t", "", "telegram bot token")
	pflag.Parse()

	app := application.ApplicationContext{
		WorkDirectory:         *workDirectory,
		UserStore:             newUserStore(*workDirectory),
		RoundStore:            newRoundStore(*workDirectory),
		LastActionStore:       newLastActionStore(*workDirectory),
		RoundParticipantStore: newRoundParticipantStore(*workDirectory),
		Bot:                   newBot(*tgToken),
	}

	log.Printf("Authorized on account %s", app.Bot.Me.Username)

	app.Bot.Handle(actions.Start, actions.OnStart(app))
	app.Bot.Handle(actions.Create, actions.OnCreate(app))
	app.Bot.Handle(actions.Join, actions.OnJoin(app))
	app.Bot.Handle(actions.PlayOut, actions.OnPlayOut(app))
	app.Bot.Handle(tb.OnText, actions.OnText(app))

	app.Bot.Start()
}
