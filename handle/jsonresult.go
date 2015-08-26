package handle

type PushResult struct {
	AppKey  string `json:"app_key"`
	Content string `json:"content"`
	UserTag string `json:"user_tag"`
	Total   int    `json:"total"`
}
type ListOnlineResult struct {
	AppKey          string   `json:"app_key"`
	TotalOnlineUser int      `json:"total_online_user"`
	OnlineUser      []string `json:"online_user"`
}
type NormalResult struct {
	Message string `json:"message"`
}

type AppResult struct {
	AppName   string `json:"app_name"`
	AppKey    string `json:"app_key"`
	RequestIP string `json:"request_ip"`
}
