package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	utils "github.com/benni347/messengerutils"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
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
	AppId          string `json:"appId"`
	AppSecret      string `json:"appSecret"`
	AppKey         string `json:"appKey"`
	ClusterId      string `json:"clusterId"`
	SupaBaseApiKey string `json:"supaBaseApiKey"`
	SupaBaseUrl    string `json:"supaBaseUrl"`
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
	a.config.SupaBaseApiKey = os.Getenv("SUPABASE_API_KEY")
	a.config.SupaBaseUrl = os.Getenv("SUPABASE_URL")
	return a.config
}

func (a *App) ValidateEmail(email string) bool {
	return utils.ValidateEmailRegex(email)
}

func (a *App) GetAppId() string {
	return a.config.AppId
}

func (a *App) GetAppSecret() string {
	return a.config.AppSecret
}

func (a *App) GetAppKey() string {
	return a.config.AppKey
}

func (a *App) GetClusterId() string {
	return a.config.ClusterId
}

func (a *App) GetSupaBaseApiKey() string {
	return a.config.SupaBaseApiKey
}

func (a *App) GetSupaBaseUrl() string {
	return a.config.SupaBaseUrl
}

type Message struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
	Time    string `json:"time"`
}

// SendMessage sends a message from a specific sender to a specified chat room.
//
// This function takes in the ID of the chat room (chatRoomId), the sender's identifier (sender), and the content of the message (message) as parameters.
// It then generates a timestamp, creates a Message object with the sender, message, and timestamp, and converts this object into a JSON format.
// Afterward, it gets the necessary Pusher credentials for this application, creates a new Pusher client using these credentials,
// and triggers a new Pusher event with the JSON message.
//
// The format of the message from the server should be: {"message": "message", "sender": "sender", "time": "time"}
//
// If there's an error while marshalling the Message object to JSON, this function prints the error and returns immediately.
//
// Note that SendMessage does not return any values.
//
// Usage:
// app := NewApp()
// app.SendMessage("chatRoomId", "senderId", "Hello, world!")
func (a *App) SendMessage(chatRoomId string, sender string, message string) {
	// The format from the server should be: {"message": "message", "time": "time"}
	currentTime := time.Now().UnixNano()
	currentTimeString := strconv.FormatInt(currentTime, 10)

	// Create a new Message object
	msg := Message{
		Message: message,
		Time:    currentTimeString,
		Sender:  sender,
	}

	// Convert Message object to JSON
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		utils.PrintError("marshalling JSON", err)
		return
	}

	// Get the Pusher credentials
	appId := a.GetAppId()
	appSecret := a.GetAppSecret()
	appKey := a.GetAppKey()
	clusterId := a.GetClusterId()

	// Create a new Pusher Client
	pusherClient := pusher.Client{
		AppID:   appId,
		Key:     appKey,
		Secret:  appSecret,
		Cluster: clusterId,
		Secure:  true,
	}

	// Create a new Pusher trigger
	pusherClient.Trigger(chatRoomId, "message", string(msgJSON))
}
