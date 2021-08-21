package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

func Run() {
	app := application{}
	app.loadConfig()
	app.init()

	ch := make(chan os.Signal)
	// run SIGINT notification
	go signal.Notify(ch, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(3)
	// run parallel processes
	go app.fireTwitchConnection(ctx, cancel, wg)
	go app.sigIntListener(ctx, cancel, wg, ch)
	go app.syncDefenderBotList(ctx, wg)

	// waiting while things are going on
	wg.Wait()
}

func (app *application) sigIntListener(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ch chan os.Signal) {
	defer wg.Done()
	// listening for 2 events.
	// 1. context was canceled by calling cancel function
	// 2. siginterrupt was received. we have to call cancel function and go away
	for {
		select {
		case <-ctx.Done():
			app.logger.Println("Context is closed")
			return
		case <-ch:
			app.logger.Println("SignInterrupt received")
			cancel()
			return
		}
	}
}
