package main

import "net"
import "fmt"
import "bufio"
//import "strings" // only needed below for sample processing

func main() {

  fmt.Println("Launching server...")
  ln, _ := net.Listen("tcp", ":8081") //на сокетах

  // accept connection on port
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)

    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    // output message received
    fmt.Print("Message Received:", string(message))
    // send new string back to client  
    conn.Write([]byte("Hello, " + message + "\n"))

}