package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/benni347/encryption"
	utils "github.com/benni347/messengerutils"
)

func Gui() {
	m := &utils.MessengerUtils{
		Verbose: true,
	}
	myApp := app.New()
	myWindow := myApp.NewWindow("Container")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}

	text1 := canvas.NewText("Hello", green)
	text2 := canvas.NewText("There", green)
	t := "There"
	tB := []byte(t)
	hash := encryption.CalculateHash(tB)
	m.PrintInfo("Hash:", hash)

	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Enter your message here")

	msgForm := &widget.Form{
		Items: []*widget.FormItem{
			{Widget: msgEntry},
		},
		OnSubmit: func() {
			msg := msgEntry.Text
			if msg == "" {
				return
			}
			msgEntry.SetText("")
			m.PrintInfo("Message:", msg)
			chatId := chatId()
			m.PrintInfo("ChatId:", chatId)
			database(msg, chatId)
		},
	}

	// content := container.NewWithoutLayout(text1, text2)
	// msgContent := container.New(layout.NewMaxLayout(), msgForm)
	content := container.NewBorder(
		container.NewAdaptiveGrid(2, text1, text2),
		nil,
		nil,
		nil,
		msgForm,
	)
	// content := container.New(layout.NewGridLayout(2), text1, text2, msgContent)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func chatId() uint64 {
	rand.Seed(time.Now().UnixNano())

	r := rand.Uint64()

	return r
}
