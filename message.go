package main

import (
	"bytes"
	"html/template"
)

var msgTemplate = `
<i>Title</i>: <b>{{.OriginalTitle}}</b>
<i>Release date</i>: <b>{{.ReleaseDate}}</b>
<i>Rating</i>: <b>{{.VoteAvg}}</b>
<i>Vote count</i>: <b>{{.VoteCount}}</b>
<i>Overview</i>: <b>{{.Overview}}</b>
{{.PosterImg}}
`

func buildMessage(m Movie) (string, error) {
	tpl, err := template.New("message_article").Parse(msgTemplate)
	if err != nil {
		return "", err
	}

	var tplBuff bytes.Buffer
	if err := tpl.Execute(&tplBuff, m); err != nil {
		return "", err
	}

	return tplBuff.String(), nil
}
