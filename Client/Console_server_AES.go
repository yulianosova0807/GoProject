package main

import (
    "net"
    "bufio"
    "fmt"
    "encoding/hex"
    "crypto/aes"
    "crypto/cipher"
    "bytes"
    "errors"
)

//Общие переменные
var (
    //Ключ для де/шифрования
    KEY, _ = hex.DecodeString("421A69BC2B99BEB97AA4BF13BE39D0344C9E31B853E646812F123DFE909F3D63")
    //Терминирующий байт
    TERM_BYTE = byte(0)
)

func main() {

  fmt.Println("Launching server...")
  ln, _ := net.Listen("tcp", ":8081") //на сокетах
  // accept connection on port
  conn, _ := ln.Accept()
	//Получаем сообщение в HEX формате
	hexRequest, _ := bufio.NewReader(conn).ReadString(TERM_BYTE)

	//Преобразовываем строку в зашифрованный массив байт
	encRequest, _ := hex.DecodeString(hexRequest[:len(hexRequest)-1])
	//Расшифровываем сообщение
	message, _ := decrypt(KEY, encRequest)

	fmt.Println("Сервер получил:", string(message))
	//Формируем ответ
	response := append(message, ", Hello"...)
	//Шифруем ответ
	encResponse, _ := encrypt(KEY, response)
	//Преобразовываем в строку
	hexResponse := hex.EncodeToString(encResponse)

	fmt.Printf("Сервер отправляет: %s\n", hexResponse)

	//Отправляем ответ обратно клиенту
	conn.Write(append([]byte(hexResponse), TERM_BYTE))
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