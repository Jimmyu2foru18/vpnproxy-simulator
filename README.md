# VPN/Proxy Simulator

A lightweight simulator for VPN/Proxy connections with visualization capabilities. This project demonstrates how proxy servers work by creating a secure tunnel between clients and target servers.

## Features

- TLS-encrypted connections between client and proxy
- Automatic certificate generation
- Tunneling to external websites
- Visual data flow representation
- Connection metrics tracking
- Configurable proxy settings

## Project Structure

```
/
├── client/
├── proxy/
├── tunnel/
├── metrics/
├── pkg/
│   ├── cert/
│   ├── logger/
│   └── visualization/
├── go.mod
├── proxy.yaml
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.20

### Installation

1. Clone the repository

```bash
git clone https://github.com/jimmymcguigan18/vpnproxy-simulator.git
cd vpnproxy-simulator
```

2. Build the proxy server

```bash
go build -o bin/proxy ./proxy
```

3. Build the client

```bash
go build -o bin/client ./client
```

## Usage

### Starting the Proxy Server

```bash
./bin/proxy

./bin/proxy --listen :9090

./bin/proxy --target example.com:80

./bin/proxy --cert mycert.pem --key mykey.pem
```

### Using the Client

```bash
./bin/client

./bin/client --proxy localhost:9090

./bin/client --target google.com

./bin/client --visual=false
```

### Understanding the Visualization

When using the client with visualization enabled (default), you'll see a real-time representation of data flowing between the client, proxy, and target server:

- **Green arrows (====>)** represent outgoing data from client to proxy
- **Cyan arrows (<====)** represent incoming data from proxy to client
- **HEX:** Shows a hexadecimal representation of the data
- **TXT:** Shows a text representation of the data (printable ASCII characters)

## How It Works

1. The proxy server listens for incoming TLS connections
2. When a client connects, the proxy reads the initial data to determine the target server
3. The proxy establishes a connection to the target server
4. Data is tunneled bidirectionally between client and target through the proxy
5. The client visualizes the data flow for educational purposes

## Configuration

The proxy server can be configured using the `proxy.yaml` file:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  max_connections: 1000

tls:
  cert_file: cert.pem
  key_file: key.pem
  min_version: 1.2

logging:
  level: info
  format: text
  file: proxy.log

metrics:
  enabled: true
  interval: 60
```