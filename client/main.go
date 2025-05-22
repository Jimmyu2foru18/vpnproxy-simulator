package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	proxyAddr = flag.String("proxy", "localhost:8080", "proxy server address")
	targetURL = flag.String("target", "example.com", "target URL to request")
	visual    = flag.Bool("visual", true, "enable visual data flow representation")
	log       = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func visualizeDataFlow(direction string, data []byte) {
	if !*visual {
		return
	}

	var arrow, color string
	if direction == "sent" {
		arrow = "====>"
		color = "\033[32m"
	} else {
		arrow = "<====" 
		color = "\033[36m" 
	}

	// Reset color
	reset := "\033[0m"

	hexView := ""
	for i, b := range data {
		if i < 16 { 
			hexView += fmt.Sprintf("%02x ", b)
		}
	}

	textView := ""
	for i, b := range data {
		if i < 32 { 
			if b >= 32 && b <= 126 { 
				textView += string(b)
			} else {
				textView += "."
			}
		}
	}

	if len(data) > 16 {
		hexView += "..."
	}
	if len(data) > 32 {
		textView += "..."
	}

	fmt.Printf("%s%s %s [%d bytes]\n", color, arrow, direction, len(data))
	fmt.Printf("%sHEX: %s\n", color, hexView)
	fmt.Printf("%sTXT: %s%s\n\n", color, textView, reset)
}

func main() {
	flag.Parse()

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", *proxyAddr, config)
	if err != nil {
		log.Fatalf("Failed to connect to proxy: %v", err)
	}
	defer conn.Close()

	log.Infof("Connected to proxy server at %s", *proxyAddr)
	log.Infof("Requesting target: %s", *targetURL)

	url := *targetURL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	request := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: VPNProxy-Simulator-Client\r\nConnection: close\r\n\r\n", *targetURL)
	visualizeDataFlow("sent", []byte(request))
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Errorf("Error sending request: %v", err)
		return
	}

	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if n > 0 {
				visualizeDataFlow("received", buffer[:n])
				os.Stdout.Write(buffer[:n])
			}
			if err != nil {
				if err != io.EOF {
					log.Errorf("Error reading from connection: %v", err)
				}
				break
			}
		}
	}()

	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buffer)
			if n > 0 {
				visualizeDataFlow("sent", buffer[:n])
				_, err = conn.Write(buffer[:n])
				if err != nil {
					log.Errorf("Error writing to connection: %v", err)
					break
				}
			}
			if err != nil {
				break
			}
		}
	}()

	for {
		time.Sleep(time.Second)
		if conn == nil {
			break
		}
	}
}