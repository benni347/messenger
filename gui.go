package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	fmt.Printf("Hash, fmt: %s", hash)

	msgEntry := widget.NewEntry()
	msgEntry.SetPlaceHolder("Enter your message here")

	msgForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Message", Widget: msgEntry},
		},
		OnSubmit: func() {
			msg := msgEntry.Text
			if msg == "" {
				return
			}
			msgEntry.SetText("")
			m.PrintInfo("Message:", msg)
		},
	}

	// content := container.NewWithoutLayout(text1, text2)
	msgContent := container.New(layout.NewMaxLayout(), msgForm)
	content := container.New(layout.NewGridLayout(2), text1, text2, msgContent)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
