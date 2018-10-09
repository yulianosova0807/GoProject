package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) { //1-ый аргумент формирует ответ сервера, 2-й запрос
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:]) //здесь начиная с 1: элемента после / пропишется в Hello
}

func main() {
	http.HandleFunc("/", handler) //обрабатывает все запросы после /,
	//например если мы напишем Http: // Localhost: 8080 / Стас в браузере,  он выведет Hello,Стас
	log.Fatal(http.ListenAndServe(":8080", nil)) //соответственно на интерфейсе :8080
}
