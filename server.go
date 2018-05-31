package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	_ "github.com/aestek/party-play/statik"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
)

var upgrader = websocket.Upgrader{}

func Serve(addr string) error {
	wsConsLock := sync.RWMutex{}
	wsCons := make(map[*websocket.Conn]struct{})

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(statikFS))

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		P.Add(q.Get("id"), &User{
			Name: q.Get("user"),
		})
	})

	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		P.Like(q.Get("id"), &User{
			Name: q.Get("user"),
		})
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print(err)
			return
		}
		wsConsLock.Lock()
		wsCons[c] = struct{}{}
		wsConsLock.Unlock()
		c.SetCloseHandler(func(code int, msg string) error {
			wsConsLock.Lock()
			delete(wsCons, c)
			wsConsLock.Unlock()
			return nil
		})

		pl, err := json.Marshal(P)
		if err != nil {
			log.Fatal(err)
		}

		err = c.WriteMessage(websocket.TextMessage, pl)
		if err != nil {
			log.Println(err)
		}
	})

	go func() {
		for range P.C {
			pl, err := json.Marshal(P)
			if err != nil {
				log.Fatal(err)
			}

			for c := range wsCons {
				err := c.WriteMessage(websocket.TextMessage, pl)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	log.Printf("listening on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
