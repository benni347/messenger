package main

import (
	"context"
	"fmt"
	"time"

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

	a.WebRTCConfig.localChannel.OnOpen(func() {
		m.PrintInfo("Local DataChannel opened")
		a.WebRTCConfig.connected = true
	})

	a.WebRTCConfig.localChannel.OnClose(func() {
		m.PrintInfo("Local DataChannel closed")
		a.WebRTCConfig.connected = false
	})

	a.WebRTCConfig.localChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		m.PrintInfo("Local DataChannel message received: " + string(msg.Data))
		a.WebRTCConfig.localMessage = string(msg.Data)
	})

	a.WebRTCConfig.remoteConnection.OnDataChannel(func(channel *webrtc.DataChannel) {
		m.PrintInfo("Remote DataChannel created")
		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			m.PrintInfo("Remote DataChannel message received: " + string(msg.Data))
			a.WebRTCConfig.remoteMessage = string(msg.Data)
		})
	})

	localOffer, err := a.WebRTCConfig.localConnection.CreateOffer(nil)
	if err != nil {
		utils.PrintError("During the local offer an error ocured", err)
		return err
	}
	m.PrintInfo("Got a local offer " + localOffer.SDP)

	err = a.WebRTCConfig.localConnection.SetLocalDescription(localOffer)
	if err != nil {
		utils.PrintError("During the local offer an error ocured", err)
		return err
	}

	err = a.WebRTCConfig.remoteConnection.SetRemoteDescription(localOffer)
	if err != nil {
		utils.PrintError("During the local offer an error ocured", err)
		return err
	}

	time.Sleep(2 * time.Second)

	remoteAnswer, err := a.WebRTCConfig.remoteConnection.CreateAnswer(nil)
	if err != nil {
		utils.PrintError("During the remote answer an error ocured", err)
		return err
	}
	m.PrintInfo("Got a remote answer " + remoteAnswer.SDP)

	err = a.WebRTCConfig.remoteConnection.SetLocalDescription(remoteAnswer)
	if err != nil {
		utils.PrintError("During the remote answer an error ocured", err)
		return err
	}

	err = a.WebRTCConfig.localConnection.SetRemoteDescription(remoteAnswer)
	if err != nil {
		utils.PrintError("During the remote answer an error ocured", err)
		return err
	}

	return nil
}

func (w *WebRTCConfig) Send(message string) error {
	err := w.localChannel.SendText(message)
	if err != nil {
		utils.PrintError("During the sending to the local channel an error ocured", err)
		return err
	}
	err = w.remoteChannel.SendText(message)
	if err != nil {
		utils.PrintError("During the sending to the remote channel an error ocured", err)
		return err
	}

	return nil
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
