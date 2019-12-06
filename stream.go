package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	hack := os.Getenv("HACK") + os.Getenv("hack")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		start := time.Now().Unix()
		if r.URL.Query().Get("hack") != "" || hack != "" {
			fmt.Printf("Got a connection (hack)\n")
			// Use a raw/hack
			hj, ok := w.(http.Hijacker)
			if !ok {
				fmt.Printf("Can't make hijack")
				http.Error(w, "Can't hijack", http.StatusInternalServerError)
				return
			}
			conn, bufrw, err := hj.Hijack()
			if err != nil {
				fmt.Printf("Can't hijack")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			buf := make([]byte, 10)
			for {
				len, err := io.ReadFull(bufrw, buf)
				if len < 10 || err != nil {
					fmt.Printf("%s Read: %s\n", time.Now().Format("15:04"), err)
					break
				}
				fmt.Printf("%s Read: %s\n", time.Now().Format("15:04"),
					string(buf[:len]))
				bufrw.Write(buf[:len])
				bufrw.Flush()
			}
		} else {
			fmt.Printf("Got a connection (ws)\n")
			// Use websockets
			upgrader := websocket.Upgrader{}

			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Printf("Upgrade failed: %s", err)
				http.Error(w, "Upgrade failed: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
			defer c.Close()
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					fmt.Printf("%s Read: %s\n", time.Now().Format("15:04"), err)
					break
				}
				fmt.Printf("%s Read: %s\n", time.Now().Format("15:04"), message)
				if err = c.WriteMessage(mt, message); err != nil {
					break
				}
			}
		}

		end := time.Now().Unix()
		fmt.Printf("\n%s Streamer duration: %d secs\n", time.Now().Format("15:04"),
			end-start)
	})

	fmt.Printf("Listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
