// 12 august 2018

// +build OMIT

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

func sendMessage(message string) string {
	client := http.Client{}
	resp, err := client.Get("http://localhost:8080/" + message)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println(err)
			return ""
		}
		bodyString := string(bodyBytes)
		return bodyString
	}
	return ""
}

var mainwin *ui.Window

func Client() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)
	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, true)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	//поле ввода текста
	group := ui.NewGroup("")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	entryForm.Append("Enter text", ui.NewEntry(), false)

	button := ui.NewButton("Connect")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		ui.MsgBox(mainwin, sendMessage("Hello"), sendMessage("I am server!"))
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	grid.Append(msggrid,
		0, 2, 2, 1,
		false, ui.AlignCenter, false, ui.AlignStart)

	return hbox

}

func setupUI() {
	mainwin = ui.NewWindow("Client_window", 300, 300, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Client", Client())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
