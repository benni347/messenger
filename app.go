package main

import (
	"context"
	"os"

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
	m.PrintInfo("The app_id is: "+a.config.AppId,
		"The app_secret is: "+a.config.AppSecret,
		"The app_key is: "+a.config.AppKey,
		"The cluster_id is: "+a.config.ClusterId)
	return a.config
}

func (a *App) ValidateEmail(email string) bool {
	return utils.ValidateEmailRegex(email)
}
