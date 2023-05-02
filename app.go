package main

import (
	"context"

	utils "github.com/benni347/messengerutils"
	webrtc "github.com/pion/webrtc/v3"
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
	verbose bool
}

type WebRTCConfig struct {
	connected        *bool
	localMessage     *string
	remoteMessage    *string
	localConnection  *webrtc.PeerConnection
	remoteConnection *webrtc.PeerConnection
	localChannel     *webrtc.DataChannel
}

type AllConfig struct {
	Config
	WebRTCConfig
	App
}

func (a *App) CreatePeerConnection() *webrtc.PeerConnection {
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError("During the PeerConnection an error ocured", err)
		panic(err)
	}

	return peerConnection
}

func (w *WebRTCConfig) Disconnect() {
	// Close the connection
	w.localConnection.Close()
	w.remoteConnection.Close()
}

func (a *AllConfig) Connect() {
	a.Config.verbose = true
	m := &utils.MessengerUtils{
		Verbose: a.Config.verbose,
	}

	m.PrintInfo("Connecting...")
	dataChannelParameters := webrtc.DataChannelParameters{
		Label:    "data",
		Ordered:  true,
		Protocol: "tcp",
	}

	localConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError("During the PeerConnection for local an error ocured", err)
		panic(err)
	}
	a.WebRTCConfig.localConnection = localConnection

	remoteConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError("During the PeerConnection for remote an error ocured", err)
		panic(err)
	}

	a.WebRTCConfig.remoteConnection = remoteConnection
}
