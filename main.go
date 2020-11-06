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

func NewStore(
	workDirectory *string) func(
	filename string,
	instanceFunction func(filename string) (interface{}, error),
) interface{} {
	return func(filename string, instanceFunction func(filename string) (interface{}, error)) interface{} {
		return utils.CheckAndReturn(instanceFunction(*workDirectory + filename))
	}
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

	newStore := NewStore(workDirectory)

	app := &application.ApplicationContext{
		*workDirectory,
		newStore(".users", func(filename string) (interface{}, error) { return store.NewUserStore(filename) }).(*store.UserStore),
		newStore(".rounds", func(filename string) (interface{}, error) { return store.NewRoundStore(filename) }).(*store.RoundStore),
		newStore(".lastactions", func(filename string) (interface{}, error) { return store.NewLastActionStore(filename) }).(*store.LastActionStore),
		newStore(".roundparticipants", func(filename string) (interface{}, error) { return store.NewRoundParticipantStore(filename) }).(*store.RoundParticipantStore),
		newStore(".roundplayouts", func(filename string) (interface{}, error) { return store.NewRoundPlayOutStore(filename) }).(*store.RoundPlayOutStore),
		newBot(*tgToken),
	}

	log.Printf("Authorized on account %s", app.Bot.Me.Username)

	app.Bot.Handle(actions.START, actions.OnStart(app))
	app.Bot.Handle(actions.CREATE, actions.OnCreate(app))
	app.Bot.Handle(actions.JOIN, actions.OnJoin(app))
	app.Bot.Handle(actions.PLAY_OUT, actions.OnPlayOut(app))
	app.Bot.Handle(tb.OnText, actions.OnText(app))

	app.Bot.Start()
}
