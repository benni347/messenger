package main

import (
	"context"
	"os"
	"time"

	utils "github.com/benni347/messengerutils"
	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx     context.Context
	verbose bool
	config  Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Config struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	AppKey    string `json:"appKey"`
	ClusterId string `json:"clusterId"`
}

// RetrieveEnvValues retrieves the values from the .env file
// and saves them in the config struct
// Returns the values of the .env file
// in the following order:
// app_id, app_secret, app_key, cluster_id
func (a *App) RetrieveEnvValues() Config {
	m := &utils.MessengerUtils{
		Verbose: a.verbose,
	}
	m.PrintInfo("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		utils.PrintError("Error loading .env file", err)
		return Config{}
	}
	a.config.AppId = os.Getenv("APP_ID")
	a.config.AppSecret = os.Getenv("SECRET")
	a.config.AppKey = os.Getenv("KEY")
	a.config.ClusterId = os.Getenv("CLUSTER")
	m.PrintInfo("The app_id is: "+a.config.AppId,
		"The app_secret is: "+a.config.AppSecret,
		"The app_key is: "+a.config.AppKey,
		"The cluster_id is: "+a.config.ClusterId)
	return a.config
}

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

func (a *App) Store(cs *Chats, chatID uint64, msg string, sender string) {
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
		Sender: sender,
	}
	chat.AddMessage(message)
}

func (a *App) Retrieve(cs *Chats, chatID uint64) []Message {
	chat := cs.FindChat(chatID)
	if chat != nil {
		return chat.Msgs
	}
	return nil
}
