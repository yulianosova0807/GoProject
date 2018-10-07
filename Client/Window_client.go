package main

import (
	"fmt"
	//"bufio"
	//"os"
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/zserge/webview"
)

const (
	windowWidth  = 480
	windowHeight = 320
)

/*<button onclick="external.invoke('connect')">Подключиться</button>
</button>
<input type="text" name="username"/>*/

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>

		<button onclick="external.invoke('Connect:'+document.getElementById('new-title').value)">
		Connect
	</button>
	<input id="new-title" type="text" />
	</body>
</html>
`

func startClient() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0") //на сокетах
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case strings.HasPrefix(data, "Connect:"):
		connect(w, strings.TrimPrefix(data, "Connect:"))
	}
}

func main() {
	url := startClient()
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "Client",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}

func connect(w webview.WebView, data string) {
	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	fmt.Fprintf(conn, data+"\n")

	message, _ := bufio.NewReader(conn).ReadString('\n')
	w.SetTitle(strings.TrimPrefix(string(message), "Connect:"))

}
