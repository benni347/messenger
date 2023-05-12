package main

import (
	"time"
)

type Message struct {
	Msg    string    `json:"msg"`
	Date   time.Time `json:"date"`
	Sender string    `json:"sender"`
}

type Chat struct {
	ChatID uint64    `json:"chatid"`
	Msgs   []Message `json:"msgs"`
}

type Chats struct {
	AllChats []Chat `json:"chats"`
}

func (c *Chat) AddMessage(msg Message) {
	c.Msgs = append(c.Msgs, msg)
}

func (cs *Chats) AddChat(chat Chat) {
	cs.AllChats = append(cs.AllChats, chat)
}

func (cs *Chats) FindChat(chatID uint64) *Chat {
	for i, chat := range cs.AllChats {
		if chat.ChatID == chatID {
			return &cs.AllChats[i]
		}
	}
	return nil
}

func store(cs *Chats, chatID uint64, msg string) {
	chat := cs.FindChat(chatID)
	if chat == nil {
		chat = &Chat{
			ChatID: chatID,
		}
		cs.AddChat(*chat)
	}
	message := Message{
		Msg:    msg,
		Date:   time.Now(),
		Sender: "you",
	}
	chat.AddMessage(message)
}

func retrieve(cs *Chats, chatID uint64) []Message {
	chat := cs.FindChat(chatID)
	if chat != nil {
		return chat.Msgs
	}
	return nil
}
