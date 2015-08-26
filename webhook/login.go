package webhook

import (
	"github.com/syhlion/gopusher/model"
	"net/http"
	"net/url"
)

type Login struct {
	appdata *model.AppData
}

func (l *Login) Pull(app_key string, token string) {
	client := &http.Client{}

	var hook_url string
	req, _ := http.NewRequest("PUT", hook_url, nil)
	req.PostForm = url.Values{"token": {token}}
	client.Do(req)
}
