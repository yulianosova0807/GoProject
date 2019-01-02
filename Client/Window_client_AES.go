package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"bufio"

	"github.com/zserge/webview"
)
import (
    "encoding/hex"
    "crypto/aes"
    "crypto/cipher"
    "bytes"
	"errors"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
)

var (
    //Ключ для де/шифрования
    KEY, _ = hex.DecodeString("421A69BC2B99BEB97AA4BF13BE39D0344C9E31B853E646812F123DFE909F3D63")
    //Терминирующий байт
    TERM_BYTE = byte(0)
)

const (
	windowWidth  = 480
	windowHeight = 320
)

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
	ln, err := net.Listen("tcp", "127.0.0.1:0")
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
		connect(w,strings.TrimPrefix(data, "Connect:"));
	}
}

func main() {
	url := startClient()
	w := webview.New(webview.Settings{
		Width:     windowWidth,
		Height:    windowHeight,
		Title:     "Client",
		Resizable: true,
		URL:       url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}

func connect(w webview.WebView, data string) {
   //Подключаемся к серверу
   conn, _ := net.Dial("tcp", "127.0.0.1:8081") 
   message := data
   //Шифруем сообщение
   encRequest, _ := encrypt(KEY, []byte(message))
   //Преобразовываем в строку
   hexRequest := hex.EncodeToString(encRequest)

   fmt.Printf("Клиент отправляет: %s\n", hexRequest)

   //Отправляем запрос на сервер
   
   conn.Write(append([]byte(hexRequest), TERM_BYTE))
   //Получаем ответ от сервера
   hexResponse, _ := bufio.NewReader(conn).ReadString(TERM_BYTE)
   //Преобразовываем строку в зашифрованный массив байт
   encResponse, _ := hex.DecodeString(hexResponse[:len(hexResponse)-1])
   
   hexAnswer := hex.EncodeToString(encResponse)
   //Расшифровываем сообщение
   messageBytes, _ := decrypt(KEY, []byte(encResponse))

   fmt.Println("Клиент получил: " + string(messageBytes))

   //узнаем время перессылки сообщения
   t := time.Now()
   fmt.Println("Время отправки сообщения:" + t.Format("2006-01-02 15:04:05"))

   w.SetTitle(strings.TrimPrefix(string(messageBytes), "Connect:"))
    daraStr :=t.Format("2006-01-02 15:04:05")
   //gпишем все данные в БД
   connStr := "user=postgres password=0000 dbname=client sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    } 
    defer db.Close()
     
    result, err :=  db.Exec("insert into Audit (client, messageToSent, HexMassageToSent, dataSent, HexMassageToGet,	MassageToGet)  values ('Client', $1, $2,$3, $4, $5)", message , hexRequest , daraStr, hexAnswer, string(messageBytes))
    if err != nil{
        panic(err)
	}
    _ = result
}

func encrypt(key []byte, message []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)

    if err != nil {
        return nil, err
    }

    b := message
    b = PKCS5Padding(b, aes.BlockSize)
    encMessage := make([]byte, len(b))
    iv := key[:aes.BlockSize]
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(encMessage, b)

    return encMessage, nil
}

func decrypt(key []byte, encMessage []byte) ([]byte, error) {
    iv := key[:aes.BlockSize]
    block, err := aes.NewCipher(key)

    if err != nil {
        return nil, err
    }

    if len(encMessage) < aes.BlockSize {
        return nil, errors.New("encMessage слишком короткий")
    }

    decrypted := make([]byte, len(encMessage))
    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(decrypted, encMessage)

    return PKCS5UnPadding(decrypted), nil
}

func PKCS5Padding(cipher []byte, blockSize int) []byte {
    padding := blockSize - len(cipher)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)

    return append(cipher, padText...)
}

func PKCS5UnPadding(src []byte) []byte {
    length := len(src)
    unPadding := int(src[length-1])

    return src[:(length - unPadding)]
}
