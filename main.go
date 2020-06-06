package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tmdb "github.com/cyruzin/golang-tmdb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
)

const appPrefix = "TMDB_BOT"

type Config struct {
	// required envs
	BotToken   string `required:"true" envconfig:"BOT_TOKEN"`
	TmdbAPIKey string `required:"true" envconfig:"TMDB_API_KEY"`

	// optional envs
	TgDebug        bool `default:"false" envconfig:"TG_DEBUG"`
	TgTimeoutInSec int  `default:"60" envconfig:"TG_TIMEOUT_IN_SEC"`
	CacheTimeInSec int  `default:"3600" envconfig:"CACHE_TIME_IN_SEC"`
}

var conf Config

func main() {
	if err := envconfig.Process(appPrefix, &conf); err != nil {
		log.Panic(err)
	}

	botAPI, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	if conf.TgDebug {
		botAPI.Debug = true
	}

	tmdbAPI, err := tmdb.Init(conf.TmdbAPIKey)
	if err != nil {
		log.Panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigCh
		log.Printf("interrupted with [%s] signal", sig)
		cancel()
	}()

	startBot(ctx, botAPI, tmdbAPI)
}
