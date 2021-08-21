package app

import (
	"github.com/gempir/go-twitch-irc/v2"
)

func (app *application) registerHandlers() {
	app.twitchClient.OnConnect(func() {
		app.logger.Println("We are connected.", app.twitchClient.IrcAddress)
		err := app.defender.SyncBanList()
		if err != nil {
			app.logger.Println("Failed initial sync of bots usernames", err)
		}
	})
	app.twitchClient.OnPrivateMessage(app.defender.Process)
	app.twitchClient.OnNoticeMessage(func(message twitch.NoticeMessage) {
		app.logger.Printf("Notice message received: %s\n", message.Message)
	})
}
