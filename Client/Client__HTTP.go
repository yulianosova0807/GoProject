// 12 august 2018

// +build OMIT

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

func sendMessage(message string) string {
	client := http.Client{}                                     //будем получать информацию для клиента по http протоколу
	resp, err := client.Get("http://localhost:8080/" + message) //метод Get - принимает адрес локалхост, и сообщение (если бы это делалось в браузере, так оно висит по умолчанию)
	if err != nil {
		fmt.Println(err) // обработка ошибок - по умолчанию всегда
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK { // статус Ок - для проверки, сработал ли GET()
		bodyBytes, err2 := ioutil.ReadAll(resp.Body) //считываем пришедшие данные ( в данном случае, это будет сообщение в окне)
		if err2 != nil {
			fmt.Println(err) // обработка ошибок - по умолчанию всегда
			return ""
		}

		bodyString := string(bodyBytes) //преобразуем байты в строку
		return bodyString               //возвражаем строку
	}
	return ""
}

var mainwin *ui.Window

func Encrypt(bodyString string) string {
	// Generate RSA Key
	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
	raulPublicKey := &raulPrivateKey.PublicKey
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, []byte(bodyString), label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(ciphertext[:])
}
func Decrypt(ciphertext string) string {
	raulPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
	label := []byte("")
	hash := sha256.New()
	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, raulPrivateKey, []byte(ciphertext), label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(plainText[:])
}

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

	group := ui.NewGroup("")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	//поле ввода текста
	textBox := ui.NewEntry() //создаем ссылку, чтобы потом можно было отсюда дергать текст
	entryForm.Append("Enter text", textBox, false)

	button := ui.NewButton("Connect")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		ui.MsgBox(mainwin, "Client_Server", sendMessage(textBox.Text())) // вызываем в МВ функцию с принятым сообщением
	})

	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button1 := ui.NewButton("Encrypt")
	entry1 := ui.NewEntry()
	entry1.SetReadOnly(true)
	button1.OnClicked(func(*ui.Button) {
		ui.MsgBox(mainwin, "Client_Server", Encrypt(textBox.Text())) // вызываем в МВ функцию с принятым сообщением
	})

	grid.Append(button1,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button2 := ui.NewButton("Encrypt")
	entry2 := ui.NewEntry()
	entry2.SetReadOnly(true)
	button2.OnClicked(func(*ui.Button) {
		ui.MsgBox(mainwin, "Client_Server", Decrypt(textBox.Text())) // вызываем в МВ функцию с принятым сообщением
	})

	grid.Append(button2,
		2, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	grid.Append(msggrid,
		0, 2, 2, 1,
		false, ui.AlignCenter, false, ui.AlignStart)

	return hbox

}

func setupUI() {
	mainwin = ui.NewWindow("", 300, 300, true)
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
