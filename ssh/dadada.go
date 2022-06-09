// https://github.com/manishmeganathan/peerchat/
package main

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	_ "embed"

)

//go:embed welcome.txt 
var welcometxtbyte []byte

//go:embed help.txt
var helptxtbyte []byte
// Represents the app version
const appversion = "v1.0.0"
var app *tview.Application

func main() {
	loadJar()
	u, p := getLogin()
	print("logging in... ")
	if !checkLogin(u,p) {
		print("Sorry! It looks like provided data was invalid.")
		os.Exit(0)
	}
	print("OK!")
	ui := NewUI()
	ui.Run()
	os.Exit(0)
}

func getLogin() (string, string) {
	app = tview.NewApplication()
	var loginusername string
	var loginpassword string
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, func(text string) {loginusername = text}).
		AddPasswordField("Password", "", 20, '*', func(text string) {loginpassword = text}).
		AddCheckbox("Are you a robot?", true, nil).
		AddButton("Login", func() {
			app.Stop()
		}).
		AddButton("Quit", func() {
			os.Exit(0)
		})
	form.SetBorder(true).SetTitle("== Login to MoKazino ==").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return loginusername, loginpassword
}
// A structure that represents the ChatRoom UI
type UI struct {
	// Represents the tview application
	TerminalApp *tview.Application

	// Represents the user message input queue
	MsgInputs chan string

	// Represents the channel of incoming messages
	Inbound chan chatmessage
	// Represents the channel of outgoing messages
	Outbound chan string
	// Represents the channel of chat log messages
	Logs chan chatlog

	// Represents the UI element with the chat messages and logs
	messageBox *tview.TextView
	// Represents the UI element for the input field
	inputBox *tview.InputField

	
}
// A structure that represents a chat message
type chatmessage struct {
	Message    string `json:"message"`
	SenderID   string `json:"senderid"`
	SenderName string `json:"sendername"`
}

// A structure that represents a chat log
type chatlog struct {
	logprefix string
	logmsg    string
}
// A structure that represents a UI command
type uicommand struct {
	cmdtype string
	cmdargs []string
}

// A constructor function that generates and
// returns a new UI for a given ChatRoom
func NewUI() *UI {
	// Create a new Tview App
	app := tview.NewApplication()

	// Initialize the command and message input channels
	msgchan := make(chan string)

	// Create a title box
	titlebox := tview.NewTextView().
		SetText(fmt.Sprintf("MoKazino CLI %s", appversion)).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignCenter)

	titlebox.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen)

	// Create a message box
	messagebox := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messagebox.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle("Logs").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite)

	// Create a usage instruction box
	//usage := tview.NewTextView().
	//	SetDynamicColors(true).
	//	SetText(`[red]/quit[green] - quit the chat | [red]/room <roomname>[green] - change chat room | [red]/user <username>[green] - change user name | [red]/clear[green] - clear the chat`)

	//usage.
	//	SetBorder(true).
	//	SetBorderColor(tcell.ColorGreen).
	//	SetTitle("Usage").
	//	SetTitleAlign(tview.AlignLeft).
	//	SetTitleColor(tcell.ColorWhite).
	//	SetBorderPadding(0, 0, 1, 0)

	// Create peer ID box
	// peerbox := tview.NewTextView()

	//peerbox.
	//	SetBorder(true).
	//	SetBorderColor(tcell.ColorGreen).
	//	SetTitle("Peers").
	//	SetTitleAlign(tview.AlignLeft).
	//	SetTitleColor(tcell.ColorWhite)

	// Create a text input box
	input := tview.NewInputField().
		SetLabel("$ ").
		SetLabelColor(tcell.ColorGreen).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle("Input").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite).
		SetBorderPadding(0, 0, 1, 0)

	// Define functionality when the input recieves a done signal (enter/tab)
	input.SetDoneFunc(func(key tcell.Key) {
		// Check if trigger was caused by a Return(Enter) press.
		if key != tcell.KeyEnter {
			return
		}

		// Read the input text
		line := input.GetText()

		// Check if there is any input text. No point printing empty messages
		if len(line) == 0 {
			return
		}

		
		// Send the message
		msgchan <- line

		// Reset the input field
		input.SetText("")
	})

	// Create a flexbox to fit all the widgets
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titlebox, 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(messagebox, 0, 1, false),
			0, 8, false).
		AddItem(input, 3, 1, true)

	// Set the flex as the app root
	app.SetRoot(flex, true)
	messagebox.Write(welcometxtbyte)
	// Create UI and return it
	return &UI{
		TerminalApp: app,
		//peerBox:     peerbox,
		messageBox:  messagebox,
		inputBox:    input,
		MsgInputs:   msgchan,
	}
}

func (ui *UI) chatMessage() {
	for {
		msg := getmessage()
		ui.messageBox.Write([]byte(msg+"\n"))
	}
}

// A method of UI that starts the UI app
func (ui *UI) Run() error {
	go ui.msgchanListen()
	go ui.chatMessage()
	defer ui.Close()
	return ui.TerminalApp.Run()
}

// A method of UI that closes the UI app
func (ui *UI) Close() {
	//ui.pscancel()
}


// A method of UI that handles a UI command
func (ui *UI) handlecommand(cmd string, args []string) {
	switch cmd {

	// Check for the quit command
	case "/quit":
		// Stop the chat UI
		ui.TerminalApp.Stop()
		return

	// Check for the clear command
	case "/clear":
		// Clear the UI message box
		ui.messageBox.Clear()

	// Check for the room change command
	case "/help":
		ui.messageBox.Write([]byte(helptxtbyte))

	case "/deposit":
		ui.messageBox.Write([]byte("You can deposit xmr to the following address:\n[orange]"+getProfile().User.MoneroAddress+"[white]\n"))
	case "/balance":
		ui.messageBox.Write([]byte("Balance: [yellow]"+XMRToDecimal(uint64(getProfile().User.Balance))+" [orange]XMR[white]\n"))
	case "/withdraw":
		if len(args) != 2 {
			ui.messageBox.Write([]byte("Invalid amount of paremeters, check /help\n"))
			return
		}
		real_amount, err := StringToXMR(args[0])
		if err != nil {
			ui.messageBox.Write([]byte(err.Error()+"\n"))
			return
		}
		ui.messageBox.Write([]byte(withdraw(real_amount, args[1])+"\n"))
	case "/dice":
		if len(args) != 3 && len(args) != 4 {
			ui.messageBox.Write([]byte(`Usage:
  /dice 0.001 49.95 lo/hi [--public] - Place a bet of 0.001 XMR on 49.95% betting low or high odds
  [--public] - send the bet output to chat, available only for bets larger than 0.01 XMR
`))
  			return
		}
		amount := args[0] // This is amount 0.001 XMR
		_chance := args[1] // Chance 49.95
		chance, err := strconv.ParseFloat(_chance, 64)
		chance = chance*100
		if err != nil {
			ui.messageBox.Write([]byte(err.Error()+"\n"))
			return
		}
		_betlohi := args[2]
		var betlo bool
		if _betlohi == "lo" {
			betlo = false
		} else if _betlohi == "hi" {
			betlo = true
		} else {
			ui.messageBox.Write([]byte(`The only valid values are lo/hi`))
			return
		}
		ui.messageBox.Write([]byte(placebet(amount, int64(chance), betlo)+"\n"))
		

	// Unsupported command
	default:
		ui.messageBox.Write([]byte(fmt.Sprintf("unsupported command - %s\n", cmd)))
	}
}

func (ui *UI) msgchanListen() {
	for {
		msg := <- ui.MsgInputs
		if strings.HasPrefix(msg, "/") {
			cmd := strings.Split(msg, " ")
			ui.messageBox.Write([]byte("$ "+msg+"\n"))
			if len(cmd) == 1 {
				ui.handlecommand(cmd[0], []string{})
			} else {
				ui.handlecommand(cmd[0], cmd[1:])
			}
		} else {
			sendmessage(msg)
		}
		// ui.messageBox.Write([]byte(msg+"\n"))
	}
}

// A method of UI that displays a message recieved from a peer
func (ui *UI) display_chatmessage(msg chatmessage) {
	prompt := fmt.Sprintf("[green]<%s>:[-]", msg.SenderName)
	fmt.Fprintf(ui.messageBox, "%s %s\n", prompt, msg.Message)
}

// A method of UI that displays a message recieved from self
func (ui *UI) display_selfmessage(msg string) {
	prompt := fmt.Sprintf("[blue]<%s>:[-]", "ui.UserName")
	fmt.Fprintf(ui.messageBox, "%s %s\n", prompt, msg)
	ui.TerminalApp.Draw()
}

// A method of UI that displays a log message
func (ui *UI) display_logmessage(log chatlog) {
	prompt := fmt.Sprintf("[yellow]<%s>:[-]", log.logprefix)
	fmt.Fprintf(ui.messageBox, "%s %s\n", prompt, log.logmsg)
	ui.TerminalApp.Draw()
}
