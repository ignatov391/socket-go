package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hubs := map[int8]*Hub{}

	fmt.Println("-------------------------------------------------------------", new(time.Time))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		fmt.Println("[mongo] err connect", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("[mongo] err ping mongo", err)
	} else {
		fmt.Println("[mongo] ping mongo SUCCESS")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
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

	err = http.ListenAndServe(*addr, handler)
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

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
