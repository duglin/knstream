package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got a connection\n")
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

		start := time.Now().Unix()
		buf := make([]byte, 10)
		for {
			len, err := io.ReadFull(bufrw, buf)
			if len < 10 {
				fmt.Printf("Read err: %s\n", err)
				break
			}
			fmt.Printf("C: %s\n", string(buf[:len]))
			bufrw.Write(buf[:len])
			bufrw.Flush()
		}
		end := time.Now().Unix()
		fmt.Printf("\nStreamer duration: %d secs\n", end-start)
	})

	fmt.Printf("Listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
