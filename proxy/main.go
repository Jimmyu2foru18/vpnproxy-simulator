package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"vpnproxy-simulator/pkg/cert"
	"vpnproxy-simulator/tunnel"
)

var (
	listenAddr = flag.String("listen", ":8080", "proxy listen address")
	targetAddr = flag.String("target", "", "target server address (default: parse from request)")
	certFile  = flag.String("cert", "cert.pem", "TLS certificate file")
	keyFile   = flag.String("key", "key.pem", "TLS key file")
	log       = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	flag.Parse()

	if err := cert.GenerateCertificate(*certFile, *keyFile); err != nil {
		log.Warnf("Failed to generate certificate: %v, will try to use existing ones", err)
	}

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", *listenAddr, config)
	if err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
	defer listener.Close()

	log.Infof("Proxy server listening on %s", *listenAddr)
	log.Infof("Certificate: %s, Key: %s", *certFile, *keyFile)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	log.Infof("New connection from %s", clientConn.RemoteAddr())

		initialData := make([]byte, 4096)
	n, err := clientConn.Read(initialData)
	if err != nil {
		log.Errorf("Failed to read initial data: %v", err)
		return
	}

	targetAddr := ""

	request := string(initialData[:n])
	if strings.HasPrefix(request, "CONNECT") {
		parts := strings.Split(request, " ")
		if len(parts) > 1 {
			targetAddr = parts[1]
		}
	} else if strings.Contains(request, "Host:") {
		lines := strings.Split(request, "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), "host:") {
				targetAddr = strings.TrimSpace(strings.TrimPrefix(line, "Host:"))
				break
			}
		}
	}

	if targetAddr == "" {
		targetAddr = "example.com:80"
	}

	log.Infof("Connecting to target: %s", target)
	targetConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Errorf("Failed to connect to target %s: %v", target, err)
		return
	}
	defer targetConn.Close()
	if _, err = targetConn.Write(initialData[:n]); err != nil {
		log.Errorf("Error writing to target: %v", err)
		return
	}

	tunnel := tunnel.NewTunnel(clientConn, targetConn)
	tunnel.Start()
	tunnel.Wait()
}
