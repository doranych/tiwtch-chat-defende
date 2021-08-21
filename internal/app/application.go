package app

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/doranych/twitch-chat-defender/internal/app/config"
	"github.com/doranych/twitch-chat-defender/internal/defender"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

// leave config short for now. Todo: defend list source config parameter, update interval etc.

type application struct {
	logger       *log.Logger
	defender     defender.Interface
	cfg          *koanf.Koanf
	twitchClient *twitch.Client
}

func (app *application) checkConfig() error {
	if app.cfg.String("twitch.username") == "" {
		return errors.New("Configuration parameter twitch.username is required. Check config.yaml")
	}
	if app.cfg.String("twitch.token") == "" {
		return errors.New("Configuration parameter twitch.token is required. Check config.yaml")
	}
	if app.cfg.String("twitch.channel") == "" {
		return errors.New("Configuration parameter twitch.channel is required. Check config.yaml")
	}
	return nil
}

func (app *application) fireTwitchConnection(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	app.twitchClient.Join(app.cfg.String("twitch.channel"))
	// listen to context cancellation
	go func() {
		for {
			select {
			case <-ctx.Done():
				err := app.twitchClient.Disconnect()
				if err != nil {
					app.logger.Println("Failed to disconnect twitch.", err)
					cancel()
					return
				}
			}
		}
	}()

	err := app.twitchClient.Connect()
	if err != nil {
		app.logger.Println("Connect error received.", err)
		cancel()
	}
}

func (app *application) init() {
	app.twitchClient = twitch.NewClient(app.cfg.String("twitch.username"), app.cfg.String("twitch.token"))
	app.logger = log.New(os.Stdout, "defender ", log.Ldate|log.Ltime)
	app.defender = defender.New(app.twitchClient, app.logger, app.cfg.String("defender.repoUrl"), app.cfg.String("defender.command"), app.cfg.String("twitch.channel"))
	app.registerHandlers()

}

func (app *application) loadConfig() {
	var err error
	app.cfg, err = config.LoadConfig()
	if err != nil {
		app.logger.Fatal("Failed to load config", err)
	}
	if err = app.checkConfig(); err != nil {
		app.logger.Fatal("Config is incorrect", err)
	}
}

func (app *application) syncDefenderBotList(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	interval := app.cfg.Duration("defender.syncInterval")
	if interval.Microseconds() != 0 { // in case it might not be specified
		t := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				err := app.defender.SyncBanList()
				if err != nil {
					app.logger.Println("Failed to sync bot list", err)
				}
			}
		}
	}
}
