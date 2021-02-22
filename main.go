package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hubs := map[int8]*Hub{}
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
		var chatID int8
		urlQuery := r.URL.Query()
		if _, ok := urlQuery["id"]; ok {
			id, _ := strconv.ParseInt(urlQuery["id"][0], 10, 8)
			chatID = int8(id)
		}

		fmt.Println("chatID", chatID)
		var hub *Hub
		if h, ok := hubs[chatID]; ok {
			hub = h
			fmt.Println("hub", hub)
		} else {
			hub = newHub()
			hubs[chatID] = hub
			go hub.run()
		}
		serveWsChat(hub, w, r)
	})

	handler := initCors(mux)

	err := http.ListenAndServe(*addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initCors(router http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	return c.Handler(router)
}
