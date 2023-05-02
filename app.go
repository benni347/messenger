package main

import (
	"context"
	"fmt"

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
	connected        bool
	localMessage     string
	remoteMessage    string
	localConnection  *webrtc.PeerConnection
	remoteConnection *webrtc.PeerConnection
	localChannel     *webrtc.DataChannel
	remoteChannel    *webrtc.DataChannel
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

func (a *AllConfig) Connect() error {
	a.Config.verbose = true
	m := &utils.MessengerUtils{
		Verbose: a.Config.verbose,
	}

	m.PrintInfo("Connecting...")
	orderedFlow := true
	dataChannelParameters := webrtc.DataChannelInit{
		Ordered: &orderedFlow,
	}

	localConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError("During the PeerConnection for local an error ocured", err)
		panic(err)
	}
	a.WebRTCConfig.localConnection = localConnection

	a.WebRTCConfig.localConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		m.PrintInfo("Local ICE Candidate: " + candidate.String())
		err := a.WebRTCConfig.remoteConnection.AddICECandidate(candidate.ToJSON())
		if err != nil {
			utils.PrintError("adding ice candidate", err)
		}
	})

	remoteConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError("During the PeerConnection for remote an error ocured", err)
		panic(err)
	}
	a.WebRTCConfig.remoteConnection = remoteConnection

	a.WebRTCConfig.remoteConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		m.PrintInfo("Remote ICE Candidate: " + candidate.String())
		err := a.WebRTCConfig.localConnection.AddICECandidate(candidate.ToJSON())
		if err != nil {
			utils.PrintError("adding ice candidate", err)
		}
	})

	localChannel, err := a.WebRTCConfig.localConnection.CreateDataChannel(
		"text",
		&dataChannelParameters,
	)
	if err != nil {
		utils.PrintError("During the DataChannel for local an error ocured", err)
		return err
	}
	a.WebRTCConfig.localChannel = localChannel
	return nil
}

func (w *WebRTCConfig) Send(message string) {
	w.localChannel.SendText(message)
	w.remoteChannel.SendText(message)
}

func (a *AllConfig) ReceiveLocalMessage() string {
	m := &utils.MessengerUtils{
		Verbose: a.Config.verbose,
	}
	m.PrintInfo(fmt.Sprintf("Recived: %s", a.WebRTCConfig.localMessage))
	return a.WebRTCConfig.localMessage
}

func (a *AllConfig) ReceiveRemoteMessage() string {
	m := &utils.MessengerUtils{
		Verbose: a.Config.verbose,
	}
	m.PrintInfo(fmt.Sprintf("Recived: %s", a.WebRTCConfig.remoteMessage))
	return a.WebRTCConfig.remoteMessage
}
