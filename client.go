package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	hack := flag.Bool("h", false, "Turn on hack path")
	port := flag.Int("p", 8080, "Port of service")
	flag.Parse()

	host := "localhost"

	if flag.NArg() > 0 {
		host = flag.Arg(0)
	}

	start := time.Now().Unix()
	var end int64

	if *hack {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, *port))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		str := fmt.Sprintf("POST /?hack=1 HTTP/1.1\r\nHost: %s\r\n\r\n", host)
		fmt.Print(str)
		if _, err = fmt.Fprintf(conn, str); err != nil {
			fmt.Printf("Error sending POST: %s\n", err)
			os.Exit(1)
		}

		go func() {
			buf := make([]byte, 10)
			for {
				if len, err := io.ReadFull(conn, buf); len < 10 || err != nil {
					fmt.Printf("%s Read: %s\n", time.Now().Format("15:04:05"), err)
					break
				}
				fmt.Printf("%s Read: %s\n", time.Now().Format("15:04:05"),
					string(buf[:10]))
			}
		}()

		buf := []byte("1234567890")
		for end == 0 {
			if _, err := conn.Write(buf); err != nil {
				fmt.Printf("\n%s Write: %s\n", time.Now().Format("15:04:05"), err)
				break
			}
			fmt.Printf("%s Write: %s\n", time.Now().Format("15:04:05"),
				string(buf[:10]))
			time.Sleep(5 * time.Second)
		}
	} else {
		query := ""
		if i := strings.Index(host, "?"); i >= 0 {
			query = host[i:]
			host = host[:i]
		}

		proto := "ws"
		if *port == 443 {
			proto = "wss"
		}

		url := fmt.Sprintf(proto+"://%s:%d/%s", host, *port, query)
		fmt.Printf("url: %s\n", url)
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Printf("dial: %s\n", err)
			os.Exit(1)
		}
		defer c.Close()

		go func() {
			for {
				if _, message, err := c.ReadMessage(); err != nil {
					fmt.Printf("%s Read: %s\n", time.Now().Format("15:04:05"), err)
					end = time.Now().Unix()
					return
				} else {
					fmt.Printf("%s Read: %s\n", time.Now().Format("15:04:05"),
						message)
				}
			}
		}()

		buf := []byte("1234567890")
		for end == 0 {
			if err := c.WriteMessage(websocket.TextMessage, buf); err != nil {
				fmt.Printf("%s Write: %s\n", time.Now().Format("15:04:05"), err)
				break
			}
			fmt.Printf("%s Write: %s\n", time.Now().Format("15:04:05"),
				string(buf))
			time.Sleep(5 * time.Second)
		}
	}

	end = time.Now().Unix()
	fmt.Printf("%s Duration: %d seconds\n", time.Now().Format("15:04:05"),
		end-start)
}
