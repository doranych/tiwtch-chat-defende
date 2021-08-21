package defender

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type (
	Interface interface {
		Process(message twitch.PrivateMessage)
		SyncBanList() error
	}
	def struct {
		sync.Mutex
		logger  *log.Logger
		cli     *twitch.Client
		ch      string
		repoUrl string
		fp      string
		bots    []string
		notBots map[string]bool
		cmd     string
	}
)

func New(twitchCli *twitch.Client, logger *log.Logger, repoUrl, command, channel string) Interface {
	if strings.HasSuffix(repoUrl, "/") {
		repoUrl = strings.TrimSuffix(repoUrl, "/")
	}
	d := &def{cli: twitchCli, logger: logger, repoUrl: repoUrl, ch: channel, cmd: command}
	d.notBots = make(map[string]bool) // create hashmap for unbanned users to prevent multiple unbans
	return d
}

func (d *def) Process(message twitch.PrivateMessage) {
	usr := message.User.Name
	d.Lock()
	ind := sort.SearchStrings(d.bots, usr)
	d.Unlock()
	// SearchStrings somehow works as binary search, but not exactly. It founds closes value, but not exactly.
	// It will bring us pain in the ass, if we will not check exact value
	if ind != len(d.bots) && d.bots[ind] == usr {
		d.cli.Say(d.ch, strings.ReplaceAll(d.cmd, "{username}", usr))
	}
}

func (d *def) SyncBanList() error {
	resp, err := http.DefaultClient.Get(d.repoUrl + "/NamelistMASTER.txt")
	if err != nil {
		return err
	}
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bots := strings.Split(string(bts), "\n")
	d.Lock()
	d.bots = bots
	d.Unlock()
	resp, err = http.DefaultClient.Get(d.repoUrl + "/False-positives.txt")
	if err != nil {
		return err
	}
	bts, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	go d.processFalsePositives(strings.Split(string(bts), "\n"))
	return nil
}

func (d *def) processFalsePositives(list []string) {
	for _, s := range list {
		if s != "" && !d.notBots[s] {
			d.cli.Say(d.ch, fmt.Sprintf("/unban %s", s))
			d.notBots[s] = true // mark this non-bot in list of unbanned
			d.logger.Printf("Unban command sent for false-positive user %s\n", s)
			time.Sleep(1 * time.Second) // to prevent rate limit overlap
		}
	}
}
