package main

import (
	"flag"
	"fmt"
	"os"

	"vpnproxy-simulator/pkg/cert"
	"github.com/sirupsen/logrus"
)

var (
	certFile = flag.String("cert", "cert.pem", "output certificate file")
	keyFile  = flag.String("key", "key.pem", "output key file")
	log      = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	flag.Parse()

		fmt.Println("VPNProxy Simulator - Certificate Generator")
	fmt.Println("===========================================\n")

		_, certErr := os.Stat(*certFile)
	_, keyErr := os.Stat(*keyFile)

	if certErr == nil || keyErr == nil {
		fmt.Printf("Warning: Certificate (%s) or key (%s) file already exists.\n", *certFile, *keyFile)
		fmt.Print("Do you want to overwrite them? (y/n): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Certificate generation cancelled.")
			return
		}
	}

		fmt.Println("Generating certificate and key...")
	err := cert.GenerateCertificate(*certFile, *keyFile)
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	fmt.Printf("\nSuccess! Certificate and key generated:\n")
	fmt.Printf("  - Certificate: %s\n", *certFile)
	fmt.Printf("  - Private Key: %s\n\n", *keyFile)
	fmt.Println("You can now start the proxy server with:")
	fmt.Printf("  go run proxy/main.go --cert %s --key %s\n", *certFile, *keyFile)
}