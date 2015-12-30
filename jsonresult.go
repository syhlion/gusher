package main

import "github.com/syhlion/gwspack"

type PushResult struct {
	AppKey  string `json:"app_key"`
	Content string `json:"content"`
	UserTag string `json:"user_tag"`
}
type ListOnlineResult struct {
	AppKey   string                      `json:"app_key"`
	Total    int                         `json:"total"`
	UserTags map[string]gwspack.UserData `json:"user_tags"`
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
	Total int             `json"total"`
	Data  []AppDataResult `json:"data"`
}
