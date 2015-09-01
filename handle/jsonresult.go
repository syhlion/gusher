package handle

import (
	"github.com/syhlion/gusher/model"
)

type PushResult struct {
	AppKey  string `json:"app_key"`
	Content string `json:"content"`
	UserTag string `json:"user_tag"`
	Total   int    `json:"total"`
}
type ListOnlineResult struct {
	AppKey   string   `json:"app_key"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Page     int      `json:"page"`
	UserTags []string `json:"user_tags"`
}
type NormalResult struct {
	Message string `json:"message"`
}

type AppResult struct {
	AppName     string `json:"app_name"`
	AppKey      string `json:"app_key"`
	ConnectHook string `json:"connect_hook"`
	RequestIP   string `json:"request_ip"`
}

type AppListResult struct {
	Limit int                   `json:"limit"`
	Page  int                   `json:"page"`
	Total int                   `json"total"`
	Data  []model.AppDataResult `json:"data"`
}
