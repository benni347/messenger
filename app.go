package main

import (
	"context"
	"os"

	utils "github.com/benni347/messengerutils"
	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx context.Context
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
	verbose   bool
	appId     string
	appSecret string
	appKey    string
	clusterId string
}

// RetrieveEnvValues retrieves the values from the .env file
// and saves them in the config struct
// Returns the values of the .env file
// in the following order:
// app_id, app_secret, app_key, cluster_id
func (c *Config) RetrieveEnvValues() (string, string, string, string) {
	m := &utils.MessengerUtils{
		Verbose: c.verbose,
	}
	err := godotenv.Load()
	if err != nil {
		utils.PrintError("Error loading .env file", err)
		return "", "", "", ""
	}
	c.appId = os.Getenv("APP_ID")
	c.appSecret = os.Getenv("APP_SECRET")
	c.appKey = os.Getenv("APP_KEY")
	c.clusterId = os.Getenv("CLUSTER")
	m.PrintInfo("The app_id is: "+c.appId, "The app_secret is: "+c.appSecret, "The app_key is: "+c.appKey, "The cluster_id is: "+c.clusterId)
	return c.appId, c.appSecret, c.appKey, c.clusterId
}
