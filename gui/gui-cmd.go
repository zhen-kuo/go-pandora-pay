package gui

import (
	"errors"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"os"
	"strconv"
)

var NotAcceptedCharacters = map[string]bool{
	"<Ctrl>":                true,
	"<Enter>":               true,
	"<MouseWheelUp>":        true,
	"<MouseWheelDown>":      true,
	"<MouseLeft>":           true,
	"<MouseRelease>":        true,
	"<Shift>":               true,
	"<Down>":                true,
	"<Up>":                  true,
	"<Left>":                true,
	"<Right>":               true,
	"<Tab>":                 true,
	"NotAcceptedCharacters": true,
}

type Command struct {
	Text     string
	Callback func(string) error
}

var commands = []Command{
	{Text: "Wallet : Decrypt"},
	{Text: "Wallet : Show Mnemnonic"},
	{Text: "Wallet : List Addresses"},
	{Text: "Wallet : Show Private Key"},
	{Text: "Wallet : Remove Address"},
	{Text: "Wallet : Create New Address"},
	{Text: "Wallet : TX: Transfer"},
	{Text: "Wallet : TX: Delegate"},
	{Text: "Wallet : TX: Withdraw"},
	{Text: "Wallet : Export JSON"},
	{Text: "Wallet : Import JSON"},
	{Text: "Exit"},
}

var cmd *widgets.List
var cmdStatus = "cmd"
var cmdInput = ""
var cmdInputCn = make(chan string)
var cmdRows []string

func CommandDefineCallback(Text string, callback func(string) error) {

	for i := range commands {
		if commands[i].Text == Text {
			commands[i].Callback = callback
			return
		}
	}

	Error(errors.New("Command " + Text + " was not found"))
}

func cmdProcess(e ui.Event) {
	switch e.ID {
	case "<C-c>":
		if cmdStatus == "read" {
			OutputRestore()
			return
		}
		os.Exit(1)
	case "<Down>":
		cmd.ScrollDown()
	case "<Up>":
		cmd.ScrollUp()
	case "<C-d>":
		cmd.ScrollHalfPageDown()
	case "<C-u>":
		cmd.ScrollHalfPageUp()
	case "<C-f>":
		cmd.ScrollPageDown()
	case "<C-b>":
		cmd.ScrollPageUp()
	case "<Home>":
		cmd.ScrollTop()
	case "<End>":
		cmd.ScrollBottom()
	case "<Enter>":

		if cmdStatus == "cmd" {
			command := commands[cmd.SelectedRow]
			cmd.SelectedRow = 0
			if command.Callback != nil {
				OutputClear()
				go func() {

					if err := command.Callback(command.Text); err != nil {
						Error(err)
					} else {
						OutputDone()
					}

				}()
			}
		} else if cmdStatus == "output done" {
			OutputRestore()
		} else if cmdStatus == "read" {
			cmdInputCn <- cmdInput
		}

	}

	if cmdStatus == "read" && !NotAcceptedCharacters[e.ID] {
		char := e.ID
		if char == "<Space>" {
			char = " "
		}
		if char == "<Backspace>" {
			char = ""
			cmdInput = cmdInput[:len(cmdInput)-1]
		}
		cmdInput = cmdInput + char
		cmd.Lock()
		cmd.Rows[len(cmd.Rows)-1] = "-> " + cmdInput
		cmd.Unlock()
	}

	// previousKey = e.ID

	ui.Render(cmd)
}

func OutputWrite(any interface{}) {
	cmd.Lock()
	cmd.Rows = append(cmd.Rows, processArgument(any))
	cmd.SelectedRow = len(cmd.Rows) - 1
	cmd.Unlock()
	ui.Render(cmd)
}

func outputRead(any interface{}) <-chan string {

	cmd.Lock()
	cmdInput = ""
	cmd.Rows = append(cmd.Rows, "")
	cmd.Rows = append(cmd.Rows, processArgument(any)+" : ")
	cmd.Rows = append(cmd.Rows, "-> ")
	cmd.SelectedRow = len(cmd.Rows) - 1
	cmdStatus = "read"
	cmd.Unlock()
	ui.Render(cmd)

	return cmdInputCn
}

func OutputReadString(any interface{}) <-chan string {
	return outputRead(any)
}

func OutputReadInt(any interface{}) <-chan int {
	r := make(chan int)

	go func() {

		for {
			str := <-outputRead(any)
			no, err := strconv.Atoi(str)
			if err != nil {
				OutputWrite("Invalid Number")
				continue
			}
			r <- no
			return
		}
	}()

	return r
}

func OutputClear() {
	cmd.Lock()
	cmd.Rows = []string{}
	cmd.Unlock()
	ui.Render(cmd)
}

func OutputDone() {
	OutputWrite("")
	OutputWrite("Press space to return...")
	cmdStatus = "output done"
}

func OutputRestore() {
	OutputClear()
	cmd.Lock()
	cmd.SelectedRow = 0
	cmd.Rows = cmdRows
	cmd.Unlock()
	ui.Render(cmd)
	cmdStatus = "cmd"
}

func cmdInit() {
	cmd = widgets.NewList()
	cmd.Title = "Commands"
	cmdRows = make([]string, len(commands))
	for i, command := range commands {
		cmdRows[i] = strconv.Itoa(i) + " " + command.Text
	}
	cmd.Rows = cmdRows
	cmd.TextStyle = ui.NewStyle(ui.ColorYellow)
	cmd.WrapText = true
}
