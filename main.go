package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/headzoo/surf.v1"
	"fmt"
	"os"
)

var app *tview.Application
var url string

func browser() {

	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		panic(err)
	}

	bodyText := bow.Dom().Find("body").Text()
	links := bow.Links()

	text := tview.NewTextView().
		SetText(bodyText).
		SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)
	
	list := tview.NewList().
		AddItem("[yellow]Exit go-net", "", 'x', func() {
			app.Stop()
			os.Exit(0)
		})
	
	input := tview.NewInputField().
		SetLabel("Search URL ").
		SetFieldWidth(184).
		SetFieldBackgroundColor(tcell.ColorYellow).
		SetFieldTextColor(tcell.ColorBlack)

	input.SetDoneFunc(func(key tcell.Key) {
		url = input.GetText()
		browser()
	})

	grid := tview.NewGrid().
		SetRows(0,2).
		SetColumns(0,37).
		AddItem(text, 0, 0, 1, 1, 0, 0, true).
		AddItem(list, 0, 1, 1, 1, 0, 0, false).
		AddItem(input, 1, 0, 1, 1, 0, 1, false)  

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF2 {
			app.SetFocus(list)
			return nil
		}
		if event.Key() == tcell.KeyF3 {
			app.SetFocus(input)
			return nil
		}
		return event
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			app.SetFocus(text)
			return nil
		}
		if event.Key() == tcell.KeyF3 {
			app.SetFocus(input)
			return nil
		}
		return event
	})

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			app.SetFocus(text)
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			app.SetFocus(list)
			return nil
		}
		return event
	})

	for _, link := range links {
		display := fmt.Sprintf("[yellow] %s", link.URL)

		list.AddItem(display, "", '+', func() {
			url = fmt.Sprintf("%s", link.URL)
			browser()
		})
	}

	text.SetBorder(true)
	list.SetBorder(true)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

}

func main() {
	app = tview.NewApplication()

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please type an URL")
		return
	}
	url = args[1]

	browser()
}
