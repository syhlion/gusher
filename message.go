package main

type Message struct {
	Action  string `json:action`
	content string `json:content`
}
