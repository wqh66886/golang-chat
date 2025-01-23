package main

import (
	"flag"
	"github.com/wqh/chat/client"
	"github.com/wqh/chat/templates"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":5678", "http service address")
	flag.Parse()

	rm := client.NewRoom()

	http.Handle("/", &templates.TemplateHandler{
		FileName: "chat.html",
	})
	http.Handle("/room", rm)

	go rm.Run()

	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
