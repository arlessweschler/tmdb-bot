package main

import (
	"bytes"
	"context"
	"log"
	"sync"
	"text/template"

	tmdb "github.com/cyruzin/golang-tmdb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func startBot(ctx context.Context, botAPI *tgbotapi.BotAPI, tmdbAPI *tmdb.Client) {
	updates, err := botAPI.GetUpdatesChan(tgbotapi.UpdateConfig{
		Limit:   0,
		Offset:  0,
		Timeout: conf.TgTimeoutInSec,
	})
	if err != nil {
		log.Printf("get update channel: %s", err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			log.Print("context interrupted")
			return
		case u, ok := <-updates:
			if !ok {
				log.Print("updates channel was closed")
				return
			}
			wg.Add(1)
			go handleMessage(ctx, &wg, u, botAPI, tmdbAPI)
		}
	}
}

func handleMessage(ctx context.Context, wg *sync.WaitGroup, update tgbotapi.Update, botAPI *tgbotapi.BotAPI, tmdbAPI *tmdb.Client) {
	defer wg.Done()
	defer log.Printf("finished to handle message %d", update.UpdateID)

	if update.InlineQuery == nil {
		log.Printf("skipped not inline message %d", update.UpdateID)
		return
	}

	log.Printf("start to handle message [%s] from [%s], message ID [%d]", update.InlineQuery.Query, update.InlineQuery.From, update.UpdateID)

	movies, nextOffset, err := searchMovies(tmdbAPI, update.InlineQuery.Query, update.InlineQuery.Offset)
	if err != nil {
		log.Print(err.Error())
		return
	}

	log.Printf("Found %d results, offset %s, for queryID %d", len(movies), update.InlineQuery.Offset, update.UpdateID)

	articles, err := buildArticles(movies)
	if err != nil {
		log.Printf("build articles: %s", err.Error())
	}

	inlineResultConfig := tgbotapi.InlineConfig{
		Results:       articles,
		InlineQueryID: update.InlineQuery.ID,
		CacheTime:     conf.CacheTimeInSec,
		IsPersonal:    false,
		NextOffset:    nextOffset,
	}

	answerResult, err := botAPI.AnswerInlineQuery(inlineResultConfig)
	if err != nil {
		log.Printf("failed to make answer query: %s", err.Error())
	}

	if !answerResult.Ok {
		log.Printf("answer query response: %s", answerResult.Result)
	}
}

func buildArticles(movies []Movie) ([]interface{}, error) {
	articles := make([]interface{}, 0, len(movies))
	for _, m := range movies {
		msg, err := buildMessage(m)
		if err != nil {
			return nil, err
		}

		articles = append(articles,
			tgbotapi.InlineQueryResultArticle{
				Type:        "article",
				ID:          m.ID,
				Title:       m.Title,
				Description: m.Overview,
				ThumbURL:    m.PosterThumb,
				InputMessageContent: tgbotapi.InputTextMessageContent{
					ParseMode: "HTML",
					Text:      msg,
				},
			},
		)
	}

	return articles, nil
}

func buildMessage(m Movie) (string, error) {
	tpl, err := template.New("message_article").ParseFiles("./templates/message_article")
	if err != nil {
		return "", err
	}

	var tplBuff bytes.Buffer
	if err := tpl.Execute(&tplBuff, m); err != nil {
		return "", err
	}

	return tplBuff.String(), nil
}
