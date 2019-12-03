package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Missing URL arg\n")
		os.Exit(1)
	}
	fmt.Printf("URL arg: %s\n", os.Args[1])
	url, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing url %q: %s\n", url, err)
		os.Exit(1)
	}
	fmt.Printf("Connecting to: %s\n", url.Host)
	conn, err := net.Dial("tcp", url.Host)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connected\n")

	str := fmt.Sprintf("POST / HTTP/1.0\r\nHost:%s\r\n\r\n", url.Host)
	fmt.Print(str)
	if _, err = fmt.Fprintf(conn, str); err != nil {
		fmt.Printf("Error sending POST: %s\n", err)
		os.Exit(1)
	}
	start := time.Now().Unix()
	var end int64

	go func() {
		buf := make([]byte, 10)
		for {
			len, err := io.ReadFull(conn, buf)
			if len < 10 {
				fmt.Printf("Read err: %s\n", err)
				end = time.Now().Unix()
				break
			}
			fmt.Printf("S: %s\n", string(buf[:10]))
		}
	}()

	buf := []byte("1234567890")
	for end == 0 {
		_, err := conn.Write(buf)
		if err != nil {
			fmt.Printf("\nWriter err: %s\n", err)
			end = time.Now().Unix()
			break
		}
		fmt.Printf("C: %s\n", string(buf[:10]))
		time.Sleep(5 * time.Second)
	}

	conn.Close()
	fmt.Printf("Duration: %d seconds\n", end-start)
}
